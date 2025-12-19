<template>
  <div class="cad-viewer" ref="containerRef">
    <div class="view-mode-switch">
      <ElButtonGroup size="small">
        <ElButton :type="viewMode === '2d' ? 'primary' : 'default'" @click="switchTo2D">
          2D
        </ElButton>
        <ElButton :type="viewMode === '3d' ? 'primary' : 'default'" @click="switchTo3D">
          3D
        </ElButton>
      </ElButtonGroup>
      <div v-if="viewMode === '3d'" class="height-control">
        <span>墙高:</span>
        <ElSlider
          v-model="wallHeight"
          :min="50"
          :max="300"
          :step="10"
          style="width: 100px"
          @change="updateWallHeight"
        />
        <span>{{ wallHeight }}cm</span>
      </div>
    </div>
    <canvas ref="canvasRef" @click="handleCanvasClick"></canvas>
    <div v-if="loading" class="loading-overlay">
      <ElIcon class="loading-icon"><Loading /></ElIcon>
      <span>正在渲染CAD图纸...</span>
    </div>
    <ElEmpty v-if="!loading && !hasData" description="暂无CAD数据" />
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted, onUnmounted, watch, computed } from 'vue'
  import { Loading } from '@element-plus/icons-vue'
  import * as THREE from 'three'
  import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls'
  import type { CadParseResult, CadEntityInfo } from '@/api/designer-cad'

  const props = defineProps<{
    parseResult?: CadParseResult | null
    visibleLayers?: string[]
  }>()

  const emit = defineEmits<{
    (e: 'entityClick', entity: CadEntityInfo): void
    (e: 'ready'): void
  }>()

  const containerRef = ref<HTMLElement>()
  const canvasRef = ref<HTMLCanvasElement>()
  const loading = ref(false)
  const viewMode = ref<'2d' | '3d'>('2d')
  const wallHeight = ref(150)

  let scene: THREE.Scene | null = null
  let camera: THREE.PerspectiveCamera | THREE.OrthographicCamera | null = null
  let renderer: THREE.WebGLRenderer | null = null
  let controls: OrbitControls | null = null
  let entityMeshes: Map<string, THREE.Object3D[]> = new Map()
  let raycaster: THREE.Raycaster | null = null
  let mouse: THREE.Vector2 | null = null
  let gridHelper: THREE.GridHelper | null = null
  let ambientLight: THREE.AmbientLight | null = null
  let directionalLight: THREE.DirectionalLight | null = null
  let cachedCenterX = 0
  let cachedCenterY = 0
  let cachedScale = 1000

  const hasData = computed(() => {
    return props.parseResult && props.parseResult.entities && props.parseResult.entities.length > 0
  })

  const colorMap: Record<number, number> = {
    1: 0xff0000,
    2: 0xffff00,
    3: 0x00ff00,
    4: 0x00ffff,
    5: 0x0000ff,
    6: 0xff00ff,
    7: 0xffffff,
    8: 0x808080,
    9: 0xc0c0c0
  }

  const wallLayerPatterns = ['墙', 'wall', 'WALL', '墙体', 'WALLS', '内墙', '外墙']

  const isWallLayer = (layerName: string): boolean => {
    return wallLayerPatterns.some((pattern) =>
      layerName.toLowerCase().includes(pattern.toLowerCase())
    )
  }

  const initThreeJS = () => {
    if (!canvasRef.value || !containerRef.value) return
    const width = containerRef.value.clientWidth
    const height = containerRef.value.clientHeight

    scene = new THREE.Scene()
    scene.background = new THREE.Color(0x1a1a2e)

    renderer = new THREE.WebGLRenderer({ canvas: canvasRef.value, antialias: true })
    renderer.setSize(width, height)
    renderer.setPixelRatio(window.devicePixelRatio)
    renderer.shadowMap.enabled = true
    renderer.shadowMap.type = THREE.PCFSoftShadowMap

    raycaster = new THREE.Raycaster()
    mouse = new THREE.Vector2()

    setup2DMode()
    animate()
    emit('ready')
  }

  const setup2DMode = () => {
    if (!containerRef.value || !renderer || !scene) return
    const width = containerRef.value.clientWidth
    const height = containerRef.value.clientHeight
    const aspect = width / height
    const frustumSize = 500

    camera = new THREE.OrthographicCamera(
      (-frustumSize * aspect) / 2,
      (frustumSize * aspect) / 2,
      frustumSize / 2,
      -frustumSize / 2,
      0.1,
      10000
    )
    camera.position.set(0, 0, 1000)
    camera.lookAt(0, 0, 0)

    if (controls) controls.dispose()
    controls = new OrbitControls(camera, renderer.domElement)
    controls.enableRotate = false
    controls.enablePan = true
    controls.enableZoom = true
    controls.screenSpacePanning = true

    if (ambientLight) scene.remove(ambientLight)
    if (directionalLight) scene.remove(directionalLight)
    updateGrid()
  }

  const setup3DMode = () => {
    if (!containerRef.value || !renderer || !scene) return
    const width = containerRef.value.clientWidth
    const height = containerRef.value.clientHeight

    camera = new THREE.PerspectiveCamera(60, width / height, 0.1, 10000)
    camera.position.set(400, 400, 400)
    camera.lookAt(0, 0, 0)

    if (controls) controls.dispose()
    controls = new OrbitControls(camera, renderer.domElement)
    controls.enableRotate = true
    controls.enablePan = true
    controls.enableZoom = true
    controls.maxPolarAngle = Math.PI / 2.1

    ambientLight = new THREE.AmbientLight(0xffffff, 0.6)
    scene.add(ambientLight)

    directionalLight = new THREE.DirectionalLight(0xffffff, 0.8)
    directionalLight.position.set(200, 500, 300)
    directionalLight.castShadow = true
    directionalLight.shadow.mapSize.width = 2048
    directionalLight.shadow.mapSize.height = 2048
    scene.add(directionalLight)

    updateGrid()
  }

  const updateGrid = () => {
    if (!scene) return
    if (gridHelper) {
      scene.remove(gridHelper)
      gridHelper.geometry.dispose()
    }

    if (viewMode.value === '2d') {
      gridHelper = new THREE.GridHelper(2000, 40, 0x444444, 0x333333)
      gridHelper.rotation.x = Math.PI / 2
    } else {
      gridHelper = new THREE.GridHelper(2000, 40, 0x444444, 0x333333)
      const groundGeometry = new THREE.PlaneGeometry(2000, 2000)
      const groundMaterial = new THREE.MeshStandardMaterial({ color: 0x1a1a2e, roughness: 0.8 })
      const ground = new THREE.Mesh(groundGeometry, groundMaterial)
      ground.rotation.x = -Math.PI / 2
      ground.position.y = -1
      ground.receiveShadow = true
      scene.add(ground)
    }
    scene.add(gridHelper)
  }

  const switchTo2D = () => {
    if (viewMode.value === '2d') return
    viewMode.value = '2d'
    setup2DMode()
    renderCadData()
  }

  const switchTo3D = () => {
    if (viewMode.value === '3d') return
    viewMode.value = '3d'
    setup3DMode()
    renderCadData()
  }

  const updateWallHeight = () => {
    if (viewMode.value === '3d') renderCadData()
  }

  const animate = () => {
    if (!renderer || !scene || !camera) return
    requestAnimationFrame(animate)
    controls?.update()
    renderer.render(scene, camera)
  }

  const renderCadData = () => {
    if (!scene || !props.parseResult) return
    loading.value = true
    clearEntities()

    const { entities, boundingBox } = props.parseResult

    // 处理boundingBox为空的情况 / Handle null boundingBox
    if (!boundingBox) {
      // 从实体中计算边界 / Calculate bounds from entities
      let minX = Infinity,
        minY = Infinity,
        maxX = -Infinity,
        maxY = -Infinity
      entities?.forEach((entity) => {
        const p = entity.properties || {}
        const coords = [p.x0, p.y0, p.x1, p.y1, p.centerX, p.centerY].filter(
          (v) => typeof v === 'number'
        ) as number[]
        coords.forEach((v, i) => {
          if (i % 2 === 0) {
            // X坐标
            minX = Math.min(minX, v)
            maxX = Math.max(maxX, v)
          } else {
            // Y坐标
            minY = Math.min(minY, v)
            maxY = Math.max(maxY, v)
          }
        })
      })
      // 如果没有有效坐标，使用默认值 / Use defaults if no valid coords
      if (!isFinite(minX)) {
        minX = 0
        maxX = 10000
        minY = 0
        maxY = 10000
      }
      cachedCenterX = (minX + maxX) / 2
      cachedCenterY = (minY + maxY) / 2
      const rangeX = maxX - minX
      const rangeY = maxY - minY
      cachedScale = Math.max(rangeX, rangeY) || 1000
    } else {
      cachedCenterX = (boundingBox.minX + boundingBox.maxX) / 2
      cachedCenterY = (boundingBox.minY + boundingBox.maxY) / 2
      const rangeX = boundingBox.maxX - boundingBox.minX
      const rangeY = boundingBox.maxY - boundingBox.minY
      cachedScale = Math.max(rangeX, rangeY) || 1000
    }

    entities.forEach((entity, index) => {
      const mesh =
        viewMode.value === '2d'
          ? createEntity2D(entity, cachedCenterX, cachedCenterY, cachedScale)
          : createEntity3D(entity, cachedCenterX, cachedCenterY, cachedScale)

      if (mesh) {
        mesh.userData = { entity, index }
        scene!.add(mesh)
        const layerMeshes = entityMeshes.get(entity.layer) || []
        layerMeshes.push(mesh)
        entityMeshes.set(entity.layer, layerMeshes)
      }
    })

    fitCameraToContent()
    loading.value = false
  }

  const createEntity2D = (
    entity: CadEntityInfo,
    centerX: number,
    centerY: number,
    scale: number
  ): THREE.Object3D | null => {
    const color = colorMap[entity.color] || 0xffffff
    const material = new THREE.LineBasicMaterial({ color })
    const p = entity.properties
    const nx = (x: number) => ((x - centerX) / scale) * 800
    const ny = (y: number) => ((y - centerY) / scale) * 800

    switch (entity.type) {
      case 'LINE': {
        const geometry = new THREE.BufferGeometry().setFromPoints([
          new THREE.Vector3(nx((p.x0 as number) || 0), ny((p.y0 as number) || 0), 0),
          new THREE.Vector3(nx((p.x1 as number) || 0), ny((p.y1 as number) || 0), 0)
        ])
        return new THREE.Line(geometry, material)
      }
      case 'CIRCLE': {
        const radius = (((p.radius as number) || 10) / scale) * 800
        const geometry = new THREE.CircleGeometry(radius, 64)
        const edges = new THREE.EdgesGeometry(geometry)
        const circle = new THREE.LineSegments(edges, material)
        circle.position.set(nx((p.x0 as number) || 0), ny((p.y0 as number) || 0), 0)
        return circle
      }
      case 'ARC': {
        const radius = (((p.radius as number) || 10) / scale) * 800
        const startAngle = (((p.startAngle as number) || 0) * Math.PI) / 180
        const endAngle = (((p.endAngle as number) || 360) * Math.PI) / 180
        const curve = new THREE.EllipseCurve(
          nx((p.x0 as number) || 0),
          ny((p.y0 as number) || 0),
          radius,
          radius,
          startAngle,
          endAngle,
          false,
          0
        )
        const points = curve.getPoints(50)
        const geometry = new THREE.BufferGeometry().setFromPoints(points)
        return new THREE.Line(geometry, material)
      }
      case 'POLYLINE':
      case 'LWPOLYLINE': {
        const points: THREE.Vector3[] = []
        let i = 0
        while (p[`x${i}`] !== undefined) {
          points.push(new THREE.Vector3(nx(p[`x${i}`] as number), ny(p[`y${i}`] as number), 0))
          i++
        }
        if (points.length < 2) return null
        const geometry = new THREE.BufferGeometry().setFromPoints(points)
        return new THREE.Line(geometry, material)
      }
      default:
        return null
    }
  }

  const createEntity3D = (
    entity: CadEntityInfo,
    centerX: number,
    centerY: number,
    scale: number
  ): THREE.Object3D | null => {
    const color = colorMap[entity.color] || 0xffffff
    const isWall = isWallLayer(entity.layer)
    const extrudeHeight = isWall ? (wallHeight.value / scale) * 800 : 5
    const p = entity.properties
    const nx = (x: number) => ((x - centerX) / scale) * 800
    const ny = (y: number) => ((y - centerY) / scale) * 800

    switch (entity.type) {
      case 'LINE': {
        const x0 = nx((p.x0 as number) || 0)
        const y0 = ny((p.y0 as number) || 0)
        const x1 = nx((p.x1 as number) || 0)
        const y1 = ny((p.y1 as number) || 0)

        if (isWall) {
          return createWall3D(x0, y0, x1, y1, extrudeHeight, color)
        } else {
          const material = new THREE.LineBasicMaterial({ color })
          const geometry = new THREE.BufferGeometry().setFromPoints([
            new THREE.Vector3(x0, 0, -y0),
            new THREE.Vector3(x1, 0, -y1)
          ])
          return new THREE.Line(geometry, material)
        }
      }
      case 'CIRCLE': {
        const x = nx((p.x0 as number) || 0)
        const y = ny((p.y0 as number) || 0)
        const radius = (((p.radius as number) || 10) / scale) * 800

        if (isWall) {
          const shape = new THREE.Shape()
          shape.absarc(0, 0, radius, 0, Math.PI * 2, false)
          const innerRadius = radius - 5
          if (innerRadius > 0) {
            const hole = new THREE.Path()
            hole.absarc(0, 0, innerRadius, 0, Math.PI * 2, true)
            shape.holes.push(hole)
          }
          const extrudeSettings = { depth: extrudeHeight, bevelEnabled: false }
          const geometry = new THREE.ExtrudeGeometry(shape, extrudeSettings)
          const mat = new THREE.MeshStandardMaterial({ color, roughness: 0.7 })
          const mesh = new THREE.Mesh(geometry, mat)
          mesh.rotation.x = -Math.PI / 2
          mesh.position.set(x, 0, -y)
          mesh.castShadow = true
          mesh.receiveShadow = true
          return mesh
        } else {
          const geometry = new THREE.RingGeometry(radius - 1, radius, 64)
          const mat = new THREE.MeshBasicMaterial({ color, side: THREE.DoubleSide })
          const mesh = new THREE.Mesh(geometry, mat)
          mesh.rotation.x = -Math.PI / 2
          mesh.position.set(x, 1, -y)
          return mesh
        }
      }
      case 'POLYLINE':
      case 'LWPOLYLINE': {
        const points2D: THREE.Vector2[] = []
        let i = 0
        while (p[`x${i}`] !== undefined) {
          points2D.push(new THREE.Vector2(nx(p[`x${i}`] as number), ny(p[`y${i}`] as number)))
          i++
        }
        if (points2D.length < 2) return null

        if (isWall) {
          const group = new THREE.Group()
          for (let j = 0; j < points2D.length - 1; j++) {
            const wall = createWall3D(
              points2D[j].x,
              points2D[j].y,
              points2D[j + 1].x,
              points2D[j + 1].y,
              extrudeHeight,
              color
            )
            if (wall) group.add(wall)
          }
          return group
        } else {
          const points3D = points2D.map((pt) => new THREE.Vector3(pt.x, 1, -pt.y))
          const geometry = new THREE.BufferGeometry().setFromPoints(points3D)
          const mat = new THREE.LineBasicMaterial({ color })
          return new THREE.Line(geometry, mat)
        }
      }
      default:
        return null
    }
  }

  const createWall3D = (
    x0: number,
    y0: number,
    x1: number,
    y1: number,
    height: number,
    color: number
  ): THREE.Mesh | null => {
    const wallThickness = 8
    const dx = x1 - x0
    const dy = y1 - y0
    const length = Math.sqrt(dx * dx + dy * dy)
    if (length < 1) return null

    const angle = Math.atan2(dy, dx)
    const geometry = new THREE.BoxGeometry(length, height, wallThickness)
    const material = new THREE.MeshStandardMaterial({ color, roughness: 0.6, metalness: 0.1 })
    const mesh = new THREE.Mesh(geometry, material)

    const centerX = (x0 + x1) / 2
    const centerY = (y0 + y1) / 2
    mesh.position.set(centerX, height / 2, -centerY)
    mesh.rotation.y = -angle
    mesh.castShadow = true
    mesh.receiveShadow = true
    return mesh
  }

  const fitCameraToContent = () => {
    if (!camera || !controls || !containerRef.value) return

    if (viewMode.value === '2d') {
      const frustumSize = 500
      const aspect = containerRef.value.clientWidth / containerRef.value.clientHeight
      const orthoCamera = camera as THREE.OrthographicCamera
      orthoCamera.left = (-frustumSize * aspect) / 2
      orthoCamera.right = (frustumSize * aspect) / 2
      orthoCamera.top = frustumSize / 2
      orthoCamera.bottom = -frustumSize / 2
      orthoCamera.updateProjectionMatrix()
      camera.position.set(0, 0, 1000)
      controls.target.set(0, 0, 0)
    } else {
      camera.position.set(500, 400, 500)
      controls.target.set(0, 0, 0)
    }
    controls.update()
  }

  const clearEntities = () => {
    if (!scene) return
    entityMeshes.forEach((meshes) => {
      meshes.forEach((mesh) => {
        scene!.remove(mesh)
        mesh.traverse((child) => {
          if (child instanceof THREE.Mesh) {
            child.geometry.dispose()
            if (Array.isArray(child.material)) {
              child.material.forEach((m) => m.dispose())
            } else {
              child.material.dispose()
            }
          } else if (child instanceof THREE.Line) {
            child.geometry.dispose()
          }
        })
      })
    })
    entityMeshes.clear()
  }

  const updateLayerVisibility = () => {
    if (!props.visibleLayers) return
    entityMeshes.forEach((meshes, layerName) => {
      const visible = props.visibleLayers!.includes(layerName)
      meshes.forEach((mesh) => {
        mesh.visible = visible
      })
    })
  }

  const handleCanvasClick = (event: MouseEvent) => {
    if (!raycaster || !mouse || !camera || !scene || !canvasRef.value) return
    const rect = canvasRef.value.getBoundingClientRect()
    mouse.x = ((event.clientX - rect.left) / rect.width) * 2 - 1
    mouse.y = -((event.clientY - rect.top) / rect.height) * 2 + 1

    raycaster.setFromCamera(mouse, camera)
    const allMeshes: THREE.Object3D[] = []
    entityMeshes.forEach((meshes) => {
      allMeshes.push(...meshes)
    })

    const intersects = raycaster.intersectObjects(allMeshes, true)
    if (intersects.length > 0) {
      let obj = intersects[0].object
      while (obj && !obj.userData.entity && obj.parent) {
        obj = obj.parent
      }
      if (obj?.userData.entity) {
        emit('entityClick', obj.userData.entity)
      }
    }
  }

  const zoomIn = () => {
    if (!camera) return
    if (viewMode.value === '2d') {
      ;(camera as THREE.OrthographicCamera).zoom *= 1.2
      camera.updateProjectionMatrix()
    } else {
      camera.position.multiplyScalar(0.8)
    }
  }

  const zoomOut = () => {
    if (!camera) return
    if (viewMode.value === '2d') {
      ;(camera as THREE.OrthographicCamera).zoom /= 1.2
      camera.updateProjectionMatrix()
    } else {
      camera.position.multiplyScalar(1.2)
    }
  }

  const resetView = () => {
    if (!camera || !controls) return
    if (viewMode.value === '2d') {
      ;(camera as THREE.OrthographicCamera).zoom = 1
      camera.position.set(0, 0, 1000)
    } else {
      camera.position.set(500, 400, 500)
    }
    controls.target.set(0, 0, 0)
    camera.updateProjectionMatrix()
    controls.update()
  }

  const handleResize = () => {
    if (!containerRef.value || !camera || !renderer) return
    const width = containerRef.value.clientWidth
    const height = containerRef.value.clientHeight

    if (viewMode.value === '2d') {
      const aspect = width / height
      const frustumSize = 500
      const orthoCamera = camera as THREE.OrthographicCamera
      orthoCamera.left = (-frustumSize * aspect) / 2
      orthoCamera.right = (frustumSize * aspect) / 2
      orthoCamera.top = frustumSize / 2
      orthoCamera.bottom = -frustumSize / 2
      orthoCamera.updateProjectionMatrix()
    } else {
      ;(camera as THREE.PerspectiveCamera).aspect = width / height
      camera.updateProjectionMatrix()
    }
    renderer.setSize(width, height)
  }

  watch(
    () => props.parseResult,
    () => {
      if (props.parseResult) {
        renderCadData()
      }
    },
    { deep: true }
  )

  watch(
    () => props.visibleLayers,
    () => {
      updateLayerVisibility()
    },
    { deep: true }
  )

  onMounted(() => {
    initThreeJS()
    window.addEventListener('resize', handleResize)
    if (props.parseResult) {
      renderCadData()
    }
  })

  onUnmounted(() => {
    window.removeEventListener('resize', handleResize)
    clearEntities()
    if (renderer) renderer.dispose()
    if (controls) controls.dispose()
  })

  defineExpose({
    zoomIn,
    zoomOut,
    resetView,
    switchTo2D,
    switchTo3D
  })
</script>

<style scoped lang="scss">
  .cad-viewer {
    position: relative;
    width: 100%;
    height: 100%;
    min-height: 400px;
    background: #1a1a2e;
    border-radius: 4px;
    overflow: hidden;

    .view-mode-switch {
      position: absolute;
      top: 12px;
      left: 12px;
      z-index: 10;
      display: flex;
      align-items: center;
      gap: 16px;
      padding: 8px 12px;
      background: rgba(0, 0, 0, 0.6);
      border-radius: 6px;
      backdrop-filter: blur(4px);

      .height-control {
        display: flex;
        align-items: center;
        gap: 8px;
        color: #fff;
        font-size: 12px;

        :deep(.el-slider) {
          .el-slider__runway {
            background: rgba(255, 255, 255, 0.2);
          }
        }
      }
    }

    canvas {
      width: 100%;
      height: 100%;
      display: block;
    }

    .loading-overlay {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      background: rgba(26, 26, 46, 0.9);
      color: #fff;
      gap: 12px;

      .loading-icon {
        font-size: 32px;
        animation: spin 1s linear infinite;
      }
    }

    :deep(.el-empty) {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);

      .el-empty__description {
        color: #999;
      }
    }
  }

  @keyframes spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }
</style>

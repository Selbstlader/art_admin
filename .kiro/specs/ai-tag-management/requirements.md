# Requirements Document

## Introduction

AI标签管理系统是一个用于统一管理 Dify 知识库配置和 AI 提示词的功能模块。通过创建可复用的"AI标签"，开发者在构建新的 AI 对话功能时无需重复编写配置代码，只需选择对应的标签即可自动关联知识库和提示词，提升开发效率并确保配置一致性。

## Glossary

- **AI_Tag**: AI标签实体，包含名称、描述、知识库ID、提示词等配置信息
- **Tag_Manager**: 标签管理系统，负责标签的增删改查操作
- **Dify_Knowledge_Base**: Dify平台的知识库，用于存储和检索领域知识
- **System_Prompt**: 系统提示词，定义AI助手的角色、行为和回答风格
- **Chat_Service**: AI对话服务，负责处理用户与AI的交互

## Requirements

### Requirement 1

**User Story:** As a developer, I want to create AI tags with associated knowledge base and prompts, so that I can reuse these configurations across different AI chat features.

#### Acceptance Criteria

1. WHEN a user submits a valid tag creation form with name, knowledge base ID, and system prompt THEN the Tag_Manager SHALL create a new AI_Tag record and return the created tag details
2. WHEN a user attempts to create a tag with a duplicate name THEN the Tag_Manager SHALL reject the creation and return a descriptive error message
3. WHEN a user attempts to create a tag with an empty name or empty system prompt THEN the Tag_Manager SHALL reject the creation and return validation error details
4. WHEN a tag is created THEN the Tag_Manager SHALL persist the tag to the database immediately

### Requirement 2

**User Story:** As a developer, I want to view and search all available AI tags, so that I can find and select the appropriate configuration for my chat feature.

#### Acceptance Criteria

1. WHEN a user requests the tag list THEN the Tag_Manager SHALL return all active AI_Tag records with pagination support
2. WHEN a user searches tags by keyword THEN the Tag_Manager SHALL return tags where name or description contains the keyword
3. WHEN displaying tag list THEN the Tag_Manager SHALL include tag name, description, knowledge base name, creation time, and status for each tag

### Requirement 3

**User Story:** As a developer, I want to update existing AI tags, so that I can modify knowledge base associations or prompts as requirements change.

#### Acceptance Criteria

1. WHEN a user submits valid tag update data THEN the Tag_Manager SHALL update the AI_Tag record and return the updated tag details
2. WHEN a user attempts to update a tag with a name that conflicts with another existing tag THEN the Tag_Manager SHALL reject the update and return a descriptive error message
3. WHEN a tag is updated THEN the Tag_Manager SHALL record the update timestamp

### Requirement 4

**User Story:** As a developer, I want to delete AI tags that are no longer needed, so that I can keep the tag list clean and manageable.

#### Acceptance Criteria

1. WHEN a user requests to delete a tag THEN the Tag_Manager SHALL perform a soft delete by setting the status to inactive
2. WHEN a user attempts to delete a tag that is currently in use by active chat sessions THEN the Tag_Manager SHALL warn the user and require confirmation before proceeding

### Requirement 5

**User Story:** As a developer, I want to select an AI tag when initiating a chat session, so that the chat automatically uses the configured knowledge base and prompts.

#### Acceptance Criteria

1. WHEN a user selects a tag and starts a chat session THEN the Chat_Service SHALL use the tag's knowledge base ID for context retrieval
2. WHEN a user selects a tag and starts a chat session THEN the Chat_Service SHALL prepend the tag's system prompt to the conversation
3. WHEN a tag has an optional custom API key configured THEN the Chat_Service SHALL use that API key instead of the default one
4. WHEN a selected tag becomes inactive during an active session THEN the Chat_Service SHALL continue using the cached configuration until the session ends

### Requirement 6

**User Story:** As a developer, I want to test an AI tag configuration before using it in production, so that I can verify the knowledge base and prompts work correctly.

#### Acceptance Criteria

1. WHEN a user clicks the test button for a tag THEN the Tag_Manager SHALL initiate a test chat session with a predefined test query
2. WHEN the test completes THEN the Tag_Manager SHALL display the AI response and indicate whether the knowledge base was successfully accessed
3. WHEN the test fails due to invalid knowledge base ID or API key THEN the Tag_Manager SHALL display a clear error message indicating the failure reason

### Requirement 7

**User Story:** As a system administrator, I want AI tag data to be serialized and deserialized correctly, so that data integrity is maintained during storage and retrieval.

#### Acceptance Criteria

1. WHEN storing AI_Tag objects to the database THEN the Tag_Manager SHALL serialize them using JSON encoding for complex fields
2. WHEN retrieving AI_Tag objects from the database THEN the Tag_Manager SHALL deserialize JSON fields back to their original structure
3. WHEN serialization or deserialization fails THEN the Tag_Manager SHALL log the error and return a descriptive error message

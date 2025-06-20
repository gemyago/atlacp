# PRD: Atlassian MCP Integration

## Introduction/Overview

The Atlassian MCP Integration is a Model Context Protocol (MCP) server that runs locally on developer laptops, enabling developers to interact with Bitbucket and Jira through natural language conversations with AI assistants. This feature eliminates the need for constant context switching between development tools by allowing developers to manage their entire development workflow—from reading Jira tickets to merging pull requests—through conversational AI within their development environment.

**Problem Statement:** Developers currently spend significant time manually switching between IDE, Jira, and Bitbucket to complete standard development workflows. This context switching disrupts focus and introduces manual overhead that can be automated.

**Goal:** Enable developers to complete 80% of their standard development workflow through natural conversation with AI, requiring only periodic confirmation and decision points.

## Goals

1. **Reduce Context Switching:** Enable developers to stay in their development environment while managing Bitbucket and Jira tasks
2. **Automate Manual Workflows:** Replace repetitive manual tasks with AI-driven automation
3. **Maintain Process Compliance:** Ensure automated workflows follow team standards and templates
4. **Provide Intelligent Assistance:** Offer smart defaults and suggestions based on project context
5. **Enable Incremental Adoption:** Allow teams to adopt features gradually without disrupting existing workflows

## User Stories

### Core Development Workflow
- **As a developer starting work on a ticket,** I want to ask "What's PROJ-123 about?" and get full context, so that I can understand requirements without opening Jira
- **As a developer ready to share my work,** I want to say "Create a PR for my feature with the usual reviewers" and have it auto-generated with proper templates, so that I don't spend time on PR setup
- **As a developer receiving feedback,** I want to update PR descriptions through conversation, so that documentation stays current without manual editing
- **As a team lead managing releases,** I want to merge PRs by saying "merge this with squash strategy" and have the system validate everything, so that our git history stays clean
- **As a developer completing features,** I want Jira tickets to automatically update when I merge PRs, so that project tracking stays accurate

### Advanced Workflows
- **As a code reviewer,** I want to review PRs following team guidelines through AI assistance, so that review quality stays consistent
- **As a developer handling complex merges,** I want AI to guide me through merge conflicts and strategy selection, so that I make informed decisions
- **As a developer maintaining branch synchronization,** I want to create branch sync PRs using a bot account and then approve them as myself, so that branch syncing follows proper review process while being automated

## Functional Requirements

### Bitbucket Integration (Priority 1)
1. **PR Creation:** The system must allow users to create pull requests through natural language, automatically populating title, description, source/destination branches
2. **PR Reading:** The system must retrieve and display pull request details, including status, reviewers, comments, and CI results
3. **PR Updates:** The system must allow users to update PR titles and descriptions while maintaining template compliance
4. **PR Approval:** The system must allow users to approve pull requests through natural language commands
5. **PR Merging:** The system must support merging PRs with strategy selection (squash, merge commit, fast-forward) and pre-merge validation
6. **Branch Sync PRs:** The system must support creating branch synchronization pull requests using a configured bot account, then allowing the developer to approve/merge them

### Jira Integration (Priority 2)
7. **Ticket Reading:** The system must retrieve and display Jira issue details by ticket number (e.g., PROJ-123) including summary, description, status, assignee, and comments
8. **Ticket Transitions:** The system must allow users to transition tickets through workflow states
9. **Label Management:** The system must support adding and removing labels from Jira issues

### Authentication & Configuration
10. **Authentication Options:** The system must authenticate with Atlassian Cloud using the simplest viable method (API tokens, basic auth, or OAuth 2.0 - with OAuth potentially deferred to later implementation phases if complexity is high)
11. **Dual Account Support:** The system must support configuration of both user credentials and bot account credentials for branch sync operations
12. **Local Configuration:** The system must store configuration locally on the user's laptop with essential settings (Atlassian workspace, default repository, bot account details)
13. **Error Handling:** The system must provide clear, actionable error messages when operations fail

### MCP Protocol Integration
14. **Tool Discovery:** The system must expose all available tools through MCP tool listing
15. **Parameter Validation:** The system must validate all input parameters and provide helpful error messages for invalid inputs
16. **Response Formatting:** The system must return structured, readable responses that AI assistants can effectively communicate to users

## Non-Goals (Out of Scope)

- **Server/Data Center Support:** MVP focuses only on Atlassian Cloud
- **Advanced Code Review Features:** Inline comments, suggestion management, detailed review workflows
- **Deployment Pipeline Integration:** Direct triggering of CI/CD pipelines
- **Custom Field Management:** Complex Jira custom field manipulation
- **Multi-Repository Orchestration:** Cross-repository dependency management
- **Team Analytics:** Productivity metrics and reporting features
- **PR Listing/Search:** Filtering and searching pull requests by various criteria
- **Jira Ticket Search:** Advanced JQL queries and ticket filtering capabilities
- **Reviewer Auto-Selection:** Automatic assignment of reviewers to pull requests

## Design Considerations

### User Experience
- **Conversational Interface:** All interactions should feel natural and require minimal technical knowledge
- **Progressive Disclosure:** Start with simple commands, reveal advanced options as needed
- **Confirmation Points:** Critical actions (merging, ticket transitions) should require user confirmation
- **Context Preservation:** Maintain conversation context across multiple tool interactions

### Error Handling
- **Clear Error Messages:** Include specific steps users can take to resolve issues
- **Permission Handling:** Clearly communicate when users lack necessary permissions

## Technical Considerations

### API Integration
- **Rate Limiting:** Handle API rate limits gracefully with appropriate user feedback
- **Authentication:** Research and implement the simplest viable authentication method for MVP
- **Error Handling:** Provide clear error messages that help users resolve issues

### Integration Requirements
- **MCP Protocol Compliance:** Must integrate seamlessly with existing MCP server architecture
- **Local Deployment:** Must run entirely on user's laptop without requiring external servers or cloud infrastructure
- **Configuration Simplicity:** Minimize configuration complexity while supporting essential functionality including dual account setup

## Success Metrics

### Functional Success
- **Tool Availability:** All defined MCP tools are discoverable and executable
- **Error Rate:** Less than 5% of user-initiated actions result in errors
- **Response Time:** 95% of operations complete within 3 seconds

### User Adoption Indicators
- **Workflow Completion:** Users can complete end-to-end development workflows (ticket → PR → merge → ticket update) through MCP tools
- **Reduced Manual Tasks:** Developers report decreased time spent on routine Bitbucket/Jira operations
- **AI Assistant Integration:** MCP tools integrate smoothly with popular AI assistants (Claude, ChatGPT, etc.)

## Open Questions

1. **Template Handling:** Research required to understand how Bitbucket repository templates work and how they should be integrated into PR creation workflow
2. **Authentication Method:** Research needed to determine the optimal authentication approach for MVP, evaluating complexity vs functionality trade-offs between API tokens, basic auth, and OAuth 2.0
3. **Bot Account Setup:** How should the system guide users through setting up and configuring bot accounts for branch sync operations? Research and user interaction is required.
4. **Merge Strategy Defaults:** How should the system choose appropriate merge strategies automatically? - no, should be explicit option to the tool
5. **Configuration Storage:** What is the optimal local storage approach for securely storing user and bot credentials on the laptop? - env variables or .env file, or json file.

---

**Document Status:** Draft  
**Target Audience:** Junior Developer  
**Implementation Phase:** MVP  
**Review Required:** Technical Lead approval needed before implementation begins 
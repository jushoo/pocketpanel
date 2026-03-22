---
name: interactive-planning
description: Interactive, conversational planning skill. Use when users want to plan a feature, refactor, or larger undertaking. No file modifications during planning - outputs a structured markdown plan to plans/ for handoff to subagents or new sessions.
---

# Interactive Planning

Interactive, conversational planning with structured markdown output for agent handoff.

## When to Use

Use this skill when the user:
- Wants to plan a new feature or refactor
- Has a complex task that needs breaking down
- Wants to think through a problem collaboratively
- Needs a plan that can be handed off to a subagent or new session

**CRITICAL**: This skill is for planning ONLY. Do not modify files, write code, or execute commands during the planning phase. Keep the conversation focused on understanding and planning.

## Planning Workflow

### Step 1: Understand the Goal

Start by understanding what the user wants to accomplish:

```
📋 Let's plan this together!

I need to understand:
1. **What** - What are you trying to build/change?
2. **Why** - What's the motivation or problem being solved?
3. **Scope** - Is this a small change or a larger undertaking?

Feel free to describe it however feels natural - I'll ask clarifying questions as needed.
```

Ask open-ended questions to understand:
- The desired outcome
- Current state vs. desired state
- Key constraints or requirements
- Success criteria

### Step 2: Explore & Clarify

**First, explore the codebase proactively.** Before asking questions that can be answered by reading code, DO read the code. This shows you understand the project and respects the user's time.

Use a questioning style that builds understanding progressively:
- Start broad, then dive into specifics
- Surface assumptions and verify them
- Identify trade-offs and help user make informed choices

#### What to Explore First

Before asking the user, try to answer these yourself by reading code:

| Question Type | How to Answer |
|---------------|---------------|
| "What does X do?" | Read the relevant file(s) |
| "How is X structured?" | Look at the directory, read main files |
| "What depends on X?" | Search for imports/references with grep |
| "How is X currently implemented?" | Read the implementation |
| "What's the tech stack for X?" | Check package.json, go.mod, imports |
| "Are there existing utilities for X?" | Search for similar functionality |
| "What's the current pattern for X?" | Read existing similar implementations |

#### When to Ask the User

Only ask the user questions that REQUIRE their input:
- **Intent/Goals** - What do they want to achieve?
- **Preferences** - Do they prefer approach A or B?
- **Constraints** - Are there business rules or deadlines?
- **Priorities** - What matters most: speed, correctness, simplicity?
- **Trade-offs** - They need to make a judgment call

#### Example Before/After

**Before (annoying):**
> "What database are you using? Let me check the codebase... Actually, let me ask you anyway to be sure. What database are you using?"

**After (smart):**
> "I see you're using SQLite with GORM. The current session table is in `internal/models/session.go`. For the new feature, I could extend this model or create a separate table. Which approach do you prefer?"

**Key principle**: If you can answer it by reading, answer it AND explain your reasoning. Then ask only what genuinely requires user input.

### Step 3: Break Down the Plan

Structure the plan as a series of logical steps. For each step:

```
### [Step N]: <Descriptive Title>
**What**: What happens in this step
**Why**: Why this step is needed
**Files/Tasks**:
- [ ] Task description
- [ ] Another task
**Acceptance**: What "done" looks like
```

### Step 4: Review & Iterate

Before writing the final plan:
```
Let me summarize what I've understood:

[Summary of the plan]

Does this capture your intent? Anything to adjust?
```

Make adjustments based on feedback until the user confirms.

## Plan Output Format

When ready, write the plan to `plans/` directory. Use this format:

### File Naming
```
plans/<YYYY-MM-DD>-<slug>-plan.md
```
Examples:
- `plans/2026-03-22-user-authentication-plan.md`
- `plans/2026-03-22-api-refactor-plan.md`

### Plan Template

```markdown
# [Title]

**Created**: YYYY-MM-DD  
**Status**: Planning | In Progress | Complete

## Overview

[2-3 sentence summary of what this plan accomplishes]

## Context

### Current State
[How things work today]

### Desired State
[How things will work after implementation]

### Motivation
[Why this change is needed]

## Scope

### In Scope
- Item 1
- Item 2

### Out of Scope
- Item 1
- Item 2

## Implementation Plan

### Step 1: [Title]
**What**: [Brief description]
**Files**: [List of files to modify or create]
**Tasks**:
- [ ] Task 1
- [ ] Task 2
**Acceptance Criteria**: [What "done" looks like]

### Step 2: [Title]
[... continue as needed]

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| [Risk description] | [High/Medium/Low] | [How to address] |

## Dependencies

- [Dependency 1] - [Why needed]
- [Dependency 2] - [Why needed]

## Testing Strategy

[How to verify the implementation works]

## Rollback Plan

[How to undo if something goes wrong]

## Checklist

- [ ] Step 1 complete
- [ ] Step 2 complete
- [ ] All tests passing
- [ ] Documentation updated
- [ ] Ready for review
```

## Tips for Agent Handoff

Write plans that are **easy to follow**:

1. **Sequential ordering** - Steps should be in execution order
2. **Explicit file paths** - Always include full paths
3. **Clear acceptance criteria** - How does the next agent know they're done?
4. **Context preservation** - Include "why" so decisions make sense later
5. **Minimal ambiguity** - Be specific about what, not just how

### Good vs Bad Instructions

**Bad** (ambiguous):
```
- Update the auth system
- Make sure it works
```

**Good** (actionable):
```
- [ ] Add session validation middleware in `apps/api/internal/middleware/session.go`
- [ ] Return 401 if session is invalid or expired
- [ ] Verify all existing auth tests pass
```

## Starting the Conversation

When the user invokes this skill with a topic, **acknowledge it directly and start engaging**. Do NOT ask them to repeat what they just said.

### If they provide a topic directly:
```
Got it - let's plan this out!

[Your topic summary/confirmation if needed, but don't just repeat them verbatim]

To fill in the details, let me first explore the codebase to understand the current implementation, then I'll ask you only what I can't determine myself.
```

**Immediately start exploring** relevant parts of the codebase. Read files, search for patterns, understand the structure. Then:
- Present what you found as facts
- Only ask about intent, preferences, or things that genuinely require user input

### If they ask to plan something without details:
Then you can ask:
```
I'll help you plan this out! 

To get started, tell me:
- What are you trying to accomplish?
- What's the current situation?
- Any constraints or requirements I should know about?

We'll work through this together step by step, and I'll write out a structured plan at the end.
```

### Key principle
The user's message IS the starting point. Build on it, don't make them repeat it. Ask specific follow-up questions to fill gaps, not generic "tell me about X" prompts.

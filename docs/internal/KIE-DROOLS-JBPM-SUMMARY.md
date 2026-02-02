# KIE/Drools/jBPM Executive Summary

**Date:** 2026-02-02
**Context:** Historical technology review during GitHub subscription cleanup

## Overview

KIE (Knowledge Is Everything) is Red Hat's umbrella project for business automation, encompassing:

| Project | Purpose |
|---------|---------|
| **Drools** | Business rules engine - declarative rule evaluation, CEP |
| **jBPM** | Business process management - workflow orchestration |
| **OptaPlanner** | Constraint satisfaction solver - scheduling, routing optimization |

All are Java/JVM-based, Apache-licensed, and now incubating under Apache as `incubator-kie-*`.

## What They Do

### Drools
- Forward/backward chaining inference engine
- Separates business logic from application code
- DMN (Decision Model and Notation) support
- Complex event processing (CEP)
- Use case: "If customer tier is gold AND order > $100 THEN apply 15% discount"

### jBPM
- BPMN 2.0 workflow execution
- Human task management
- Process versioning and migration
- Use case: Approval workflows, onboarding processes

### OptaPlanner
- Constraint satisfaction/optimization
- Algorithms: local search, simulated annealing, tabu search
- Use case: Employee scheduling, vehicle routing, resource allocation

## Historical Connection

These subscriptions date from ~2013-2016 when Jeff worked at Pentaho/Hitachi Vantara, which had integrations with Red Hat's business automation stack. The KIE workbench (kie-wb) provided a unified UI for rules and process authoring.

## Relevance to Current Work

| Current Focus | KIE Overlap | Verdict |
|---------------|-------------|---------|
| Agent orchestration (project-conductor) | jBPM does workflow orchestration | Different paradigm - LLM vs declarative |
| Capability catalog | Drools has rule-based decisions | Could inform design, but overkill |
| Terminal UIs (C/Go) | KIE is Java-only | No tech stack overlap |
| Claude Code skills | OptaPlanner does optimization | Interesting but not needed now |

**Conclusion:** Conceptual similarities exist (orchestration, rules, optimization) but implementation approaches differ fundamentally. Current work uses LLM-based reasoning rather than declarative rule engines. No active integration planned.

## Action Taken

- Unsubscribed from 47 upstream KIE/Drools/jBPM/Pentaho repos
- Unsubscribed from 17 personal forks of same
- Retained knowledge of concepts for potential future reference

## If Revisiting Later

OptaPlanner's constraint solving could theoretically complement LLM-based agents for:
- Multi-agent task scheduling with hard constraints
- Resource-bounded planning problems
- Deterministic optimization where LLM reasoning is insufficient

Drools' rule engine could provide:
- Guardrails/policies for agent behavior
- Deterministic decision paths for compliance scenarios
- Audit trails for rule-based decisions

These would require Java interop or rewriting concepts in Go/C.

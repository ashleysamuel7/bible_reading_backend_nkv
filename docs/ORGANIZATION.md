# Documentation Organization

This document describes how documentation is organized in the backend.

## Organization Summary

All markdown documentation files have been moved from the root directory into the `docs/` folder and organized by category.

## Directory Structure

```
docs/
├── api/                          # API Documentation
│   ├── README.md
│   └── API_REFERENCE.md
│
├── details/                      # Feature Details
│   ├── arc.md
│   ├── ux.md
│   ├── ux_review.md
│   └── voice_feature.md
│
├── guides/                        # Getting Started Guides
│   ├── README.md
│   └── QUICK_START.md
│
├── implementation/                # Implementation Documentation
│   ├── README.md
│   ├── README_USER_MANAGEMENT.md
│   └── USER_MANAGEMENT_IMPLEMENTATION.md
│
├── reports/                       # Test Reports & Analysis
│   ├── README.md
│   ├── IMPLEMENTATION_TEST_REPORT.md
│   └── TEST_FAILURE_ANALYSIS.md
│
├── rules/                         # Development Rules
│
├── tests/                         # Testing Documentation
│   ├── README.md
│   ├── generateTestCase.md
│   ├── run_tests.sh
│   └── test_cases.json
│
├── PROJECT_STRUCTURE.md          # Project Structure Overview
└── README.md                      # Main Documentation Index
```

## File Categories

### API Documentation (`api/`)
- **API_REFERENCE.md** - Complete API reference with all endpoints, authentication, request/response formats

### Guides (`guides/`)
- **QUICK_START.md** - Quick start guide for setting up and running the backend

### Implementation (`implementation/`)
- **USER_MANAGEMENT_IMPLEMENTATION.md** - Detailed implementation guide for user management
- **README_USER_MANAGEMENT.md** - User management feature overview

### Reports (`reports/`)
- **IMPLEMENTATION_TEST_REPORT.md** - Test report for implementation verification
- **TEST_FAILURE_ANALYSIS.md** - Analysis of test failures

### Details (`details/`)
- **arc.md** - Architecture documentation
- **ux.md** - User experience considerations
- **ux_review.md** - UX review documentation
- **voice_feature.md** - Voice feature documentation

### Tests (`tests/`)
- Testing documentation, test case generation, and test execution scripts

### Root Level
- **PROJECT_STRUCTURE.md** - Project structure overview

## Migration Notes

The following files were moved from the root directory:
- `API_REFERENCE.md` → `docs/api/`
- `QUICK_START.md` → `docs/guides/`
- `PROJECT_STRUCTURE.md` → `docs/`
- `IMPLEMENTATION_TEST_REPORT.md` → `docs/reports/`
- `TEST_FAILURE_ANALYSIS.md` → `docs/reports/`
- `USER_MANAGEMENT_IMPLEMENTATION.md` → `docs/implementation/`
- `README_USER_MANAGEMENT.md` → `docs/implementation/`

## Related Documentation

- Test execution reports: `../../test/reports/`
- Full-stack documentation: `../../common/docs/`
- Frontend documentation: `../../bible_reading_frontend_go/docs/`


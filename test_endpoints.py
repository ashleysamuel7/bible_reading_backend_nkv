#!/usr/bin/env python3
"""
Comprehensive endpoint testing script for Bible Reading Backend
Tests all endpoints and generates a detailed report
"""

import requests
import json
import time
from datetime import datetime
from typing import Dict, List, Any
import sys
import os

# Configuration
BASE_URL = "http://localhost:8000"
TIMEOUT = 30

# Test paths configuration
TEST_DIR = os.path.join(os.path.dirname(__file__), "test")
REPORTS_DIR = os.path.join(TEST_DIR, "reports")
FIXTURES_DIR = os.path.join(TEST_DIR, "fixtures")

# Ensure directories exist
os.makedirs(REPORTS_DIR, exist_ok=True)
os.makedirs(FIXTURES_DIR, exist_ok=True)

# Colors for terminal output
class Colors:
    GREEN = '\033[92m'
    RED = '\033[91m'
    YELLOW = '\033[93m'
    BLUE = '\033[94m'
    RESET = '\033[0m'
    BOLD = '\033[1m'

class EndpointTestResult:
    def __init__(self, name: str, method: str, url: str):
        self.name = name
        self.method = method
        self.url = url
        self.request_body = None
        self.request_headers = {}
        self.response_status = None
        self.response_headers = {}
        self.response_body = None
        self.error = None
        self.duration = 0
        self.success = False

def test_endpoint(name: str, method: str, url: str, headers: Dict = None, 
                  body: Dict = None, params: Dict = None) -> EndpointTestResult:
    """Test a single endpoint"""
    result = EndpointTestResult(name, method, url)
    
    full_url = f"{BASE_URL}{url}"
    
    try:
        start_time = time.time()
        
        if method == "GET":
            response = requests.get(full_url, headers=headers, params=params, timeout=TIMEOUT)
        elif method == "POST":
            response = requests.post(full_url, headers=headers, json=body, timeout=TIMEOUT)
        elif method == "PUT":
            response = requests.put(full_url, headers=headers, json=body, timeout=TIMEOUT)
        elif method == "DELETE":
            response = requests.delete(full_url, headers=headers, timeout=TIMEOUT)
        else:
            result.error = f"Unsupported method: {method}"
            return result
        
        result.duration = time.time() - start_time
        result.response_status = response.status_code
        result.response_headers = dict(response.headers)
        
        try:
            result.response_body = response.json()
        except:
            result.response_body = response.text[:500]  # Limit text response
        
        result.request_body = body
        result.request_headers = headers or {}
        result.success = 200 <= response.status_code < 300
        
    except requests.exceptions.ConnectionError:
        result.error = "Connection refused - Server may not be running"
    except requests.exceptions.Timeout:
        result.error = f"Request timed out after {TIMEOUT} seconds"
    except Exception as e:
        result.error = str(e)
    
    return result

def print_test_result(result: EndpointTestResult):
    """Print test result to console"""
    status_color = Colors.GREEN if result.success else Colors.RED
    status_text = "✓ PASS" if result.success else "✗ FAIL"
    
    print(f"\n{Colors.BOLD}{result.name}{Colors.RESET}")
    print(f"  {status_color}{status_text}{Colors.RESET} - {result.method} {result.url}")
    
    if result.error:
        print(f"  {Colors.RED}Error: {result.error}{Colors.RESET}")
    else:
        print(f"  Status: {result.response_status}")
        print(f"  Duration: {result.duration:.3f}s")
        
        if result.request_body:
            print(f"  Request Body: {json.dumps(result.request_body, indent=2)}")
        
        if result.response_body:
            if isinstance(result.response_body, dict):
                # Show first few keys for large responses
                if len(result.response_body) > 5:
                    preview = {k: result.response_body[k] for k in list(result.response_body.keys())[:3]}
                    print(f"  Response Preview: {json.dumps(preview, indent=2)}...")
                    print(f"  Response Size: {len(json.dumps(result.response_body))} characters")
                else:
                    print(f"  Response: {json.dumps(result.response_body, indent=2)}")
            else:
                print(f"  Response: {str(result.response_body)[:200]}...")

def generate_report(results: List[EndpointTestResult]) -> str:
    """Generate a detailed markdown report"""
    report = []
    report.append("# Bible Reading Backend - API Endpoint Test Report")
    report.append(f"\n**Generated:** {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    report.append(f"**Base URL:** {BASE_URL}\n")
    
    # Summary
    total = len(results)
    passed = sum(1 for r in results if r.success)
    failed = total - passed
    
    report.append("## 📊 Summary")
    report.append(f"- **Total Endpoints:** {total}")
    report.append(f"- **Passed:** {passed} ✓")
    report.append(f"- **Failed:** {failed} ✗")
    report.append(f"- **Success Rate:** {(passed/total*100):.1f}%\n")
    
    # Detailed Results
    report.append("## 📋 Detailed Test Results\n")
    
    for i, result in enumerate(results, 1):
        status_icon = "✓" if result.success else "✗"
        status_text = "PASS" if result.success else "FAIL"
        
        report.append(f"### {i}. {result.name} {status_icon}")
        report.append(f"\n**Status:** `{status_text}`")
        report.append(f"**Method:** `{result.method}`")
        report.append(f"**Endpoint:** `{result.url}`")
        report.append(f"**HTTP Status:** `{result.response_status or 'N/A'}`")
        report.append(f"**Duration:** `{result.duration:.3f}s`")
        
        if result.error:
            report.append(f"\n**Error:** `{result.error}`")
        
        if result.request_headers:
            report.append(f"\n**Request Headers:**")
            report.append("```json")
            report.append(json.dumps(result.request_headers, indent=2))
            report.append("```")
        
        if result.request_body:
            report.append(f"\n**Request Body:**")
            report.append("```json")
            report.append(json.dumps(result.request_body, indent=2))
            report.append("```")
        
        if result.response_headers:
            report.append(f"\n**Response Headers:**")
            report.append("```json")
            report.append(json.dumps(dict(list(result.response_headers.items())[:5]), indent=2))
            report.append("```")
        
        if result.response_body:
            report.append(f"\n**Response Body:**")
            report.append("```json")
            
            if isinstance(result.response_body, dict):
                # For large responses, show structure + sample
                if len(str(result.response_body)) > 2000:
                    # Show first few items
                    preview = {}
                    if isinstance(result.response_body, list) and len(result.response_body) > 0:
                        preview = result.response_body[:2]
                        report.append(f"// Array with {len(result.response_body)} items, showing first 2:")
                    else:
                        keys = list(result.response_body.keys())[:5]
                        preview = {k: result.response_body[k] for k in keys}
                        report.append(f"// Response has {len(result.response_body)} keys, showing first 5:")
                    
                    report.append(json.dumps(preview, indent=2))
                    report.append(f"// ... (truncated, total size: {len(json.dumps(result.response_body))} chars)")
                else:
                    report.append(json.dumps(result.response_body, indent=2))
            else:
                report.append(str(result.response_body)[:1000])
            
            report.append("```")
        
        report.append("\n---\n")
    
    # Failed Tests Summary
    failed_tests = [r for r in results if not r.success]
    if failed_tests:
        report.append("## ❌ Failed Tests Summary\n")
        for result in failed_tests:
            report.append(f"- **{result.name}** ({result.method} {result.url})")
            report.append(f"  - Status: {result.response_status or 'N/A'}")
            report.append(f"  - Error: {result.error or 'N/A'}\n")
    
    return "\n".join(report)

def main():
    """Run all endpoint tests"""
    print(f"{Colors.BOLD}{Colors.BLUE}")
    print("=" * 60)
    print("Bible Reading Backend - API Endpoint Testing")
    print("=" * 60)
    print(f"{Colors.RESET}")
    
    results = []
    
    # Test 1: Health Check - Readiness
    print(f"\n{Colors.YELLOW}Testing Health Endpoints...{Colors.RESET}")
    result = test_endpoint(
        "Health Check - Readiness",
        "GET",
        "/readiness"
    )
    results.append(result)
    print_test_result(result)
    
    # Test 2: Health Check - Liveness
    result = test_endpoint(
        "Health Check - Liveness",
        "GET",
        "/liveness"
    )
    results.append(result)
    print_test_result(result)
    
    # Test 3: Get All Verses
    print(f"\n{Colors.YELLOW}Testing NIV API Endpoints...{Colors.RESET}")
    result = test_endpoint(
        "Get All Verses",
        "GET",
        "/api/niv/verses"
    )
    results.append(result)
    print_test_result(result)
    
    # Test 4: Get All Books
    result = test_endpoint(
        "Get All Books",
        "GET",
        "/api/niv/books"
    )
    results.append(result)
    print_test_result(result)
    
    # Get a book name from the books endpoint for subsequent tests
    book_name = None
    if results[-1].success and results[-1].response_body:
        if isinstance(results[-1].response_body, list) and len(results[-1].response_body) > 0:
            book_name = results[-1].response_body[0].get('book', 'Genesis')
        elif isinstance(results[-1].response_body, dict):
            book_name = results[-1].response_body.get('book', 'Genesis')
    
    if not book_name:
        book_name = "Genesis"  # Default fallback
    
    # Test 5: Get Chapters for a Book
    result = test_endpoint(
        f"Get Chapters for Book: {book_name}",
        "GET",
        f"/api/niv/chapters/{book_name}"
    )
    results.append(result)
    print_test_result(result)
    
    # Test 6: Get Verses by Chapter (valid)
    chapter = 1
    result = test_endpoint(
        f"Get Verses by Chapter - {book_name} Chapter {chapter}",
        "GET",
        f"/api/niv/{book_name}/{chapter}/verses"
    )
    results.append(result)
    print_test_result(result)
    
    # Test 7: Get Verses by Chapter (invalid chapter)
    result = test_endpoint(
        f"Get Verses by Chapter - Invalid Chapter (999)",
        "GET",
        f"/api/niv/{book_name}/999/verses"
    )
    results.append(result)
    print_test_result(result)
    
    # Test 8: Explain Verse (with valid data)
    explain_body = {
        "book": book_name,
        "chapter": 1,
        "start_verse": 1,
        "end_verse": 3,
        "age": 25,
        "belief": 3
    }
    result = test_endpoint(
        "Explain Verse - Valid Request",
        "POST",
        "/api/niv/explain",
        headers={"Content-Type": "application/json"},
        body=explain_body
    )
    results.append(result)
    print_test_result(result)
    
    # Test 9: Explain Verse (invalid request - missing fields)
    explain_body_invalid = {
        "book": book_name,
        "chapter": 1
        # Missing required fields
    }
    result = test_endpoint(
        "Explain Verse - Invalid Request (Missing Fields)",
        "POST",
        "/api/niv/explain",
        headers={"Content-Type": "application/json"},
        body=explain_body_invalid
    )
    results.append(result)
    print_test_result(result)
    
    # Test 10: Explain Verse (with default age/belief)
    explain_body_defaults = {
        "book": book_name,
        "chapter": 1,
        "start_verse": 1,
        "end_verse": 3
        # age and belief will use defaults (25, 3)
    }
    result = test_endpoint(
        "Explain Verse - With Default Values",
        "POST",
        "/api/niv/explain",
        headers={"Content-Type": "application/json"},
        body=explain_body_defaults
    )
    results.append(result)
    print_test_result(result)
    
    # Test 11-30: User Management Endpoints
    print(f"\n{Colors.YELLOW}Testing User Management Endpoints...{Colors.RESET}")
    
    # Test 11: Register User (Valid)
    test_email = f"testuser_{int(time.time())}@example.com"
    register_body = {
        "first_name": "Test",
        "last_name": "User",
        "email": test_email,
        "password": "testpass123",
        "age": 25,
        "belif_rating": 3
    }
    result = test_endpoint(
        "Register User - Valid",
        "POST",
        "/api/register/",
        headers={"Content-Type": "application/json"},
        body=register_body
    )
    results.append(result)
    print_test_result(result)
    
    # Extract token from registration response
    auth_token = None
    user_id = None
    if result.success and result.response_body and isinstance(result.response_body, dict):
        auth_token = result.response_body.get("access")
        user_data = result.response_body.get("user", {})
        user_id = user_data.get("id")
    
    # Test 12: Register User (Duplicate Email)
    result = test_endpoint(
        "Register User - Duplicate Email",
        "POST",
        "/api/register/",
        headers={"Content-Type": "application/json"},
        body=register_body
    )
    results.append(result)
    print_test_result(result)
    
    # Test 13: Register User (Invalid - Missing Fields)
    result = test_endpoint(
        "Register User - Missing Required Fields",
        "POST",
        "/api/register/",
        headers={"Content-Type": "application/json"},
        body={"email": "test@example.com"}  # Missing required fields
    )
    results.append(result)
    print_test_result(result)
    
    # Test 14: Register User (Invalid - Short Password)
    result = test_endpoint(
        "Register User - Short Password",
        "POST",
        "/api/register/",
        headers={"Content-Type": "application/json"},
        body={
            "first_name": "Test",
            "last_name": "User",
            "email": f"test2_{int(time.time())}@example.com",
            "password": "123",  # Too short
            "age": 25,
            "belif_rating": 3
        }
    )
    results.append(result)
    print_test_result(result)
    
    # Test 15: Login User (Valid)
    login_body = {
        "email": test_email,
        "password": "testpass123"
    }
    result = test_endpoint(
        "Login User - Valid Credentials",
        "POST",
        "/api/login/",
        headers={"Content-Type": "application/json"},
        body=login_body
    )
    results.append(result)
    print_test_result(result)
    
    # Update token from login if registration didn't provide it
    if result.success and result.response_body and isinstance(result.response_body, dict):
        auth_token = result.response_body.get("access") or auth_token
    
    # Test 16: Login User (Invalid Credentials)
    result = test_endpoint(
        "Login User - Invalid Password",
        "POST",
        "/api/login/",
        headers={"Content-Type": "application/json"},
        body={
            "email": test_email,
            "password": "wrongpassword"
        }
    )
    results.append(result)
    print_test_result(result)
    
    # Test 17: Login User (Non-existent Email)
    result = test_endpoint(
        "Login User - Non-existent Email",
        "POST",
        "/api/login/",
        headers={"Content-Type": "application/json"},
        body={
            "email": "nonexistent@example.com",
            "password": "password123"
        }
    )
    results.append(result)
    print_test_result(result)
    
    # If we have a token, test protected endpoints
    if auth_token:
        auth_headers = {
            "Authorization": f"Bearer {auth_token}",
            "Content-Type": "application/json"
        }
        
        # Test 18: Get Current User Profile
        result = test_endpoint(
            "Get Current User Profile",
            "GET",
            "/api/users/me",
            headers={"Authorization": f"Bearer {auth_token}"}
        )
        results.append(result)
        print_test_result(result)
        
        # Test 19: Update Current User Profile
        result = test_endpoint(
            "Update Current User Profile",
            "PUT",
            "/api/users/me",
            headers=auth_headers,
            body={
                "first_name": "Updated",
                "age": 30,
                "belif_rating": 4
            }
        )
        results.append(result)
        print_test_result(result)
        
        # Test 20: Update Current User Profile (Invalid Age)
        result = test_endpoint(
            "Update Current User Profile - Invalid Age",
            "PUT",
            "/api/users/me",
            headers=auth_headers,
            body={
                "age": 200  # Invalid age
            }
        )
        results.append(result)
        print_test_result(result)
        
        # Test 21: Add Favorite Verse
        result = test_endpoint(
            "Add Favorite Verse",
            "POST",
            "/api/users/me/favorites",
            headers=auth_headers,
            body={
                "book_id": 1,
                "chapter": 1,
                "verse": 1
            }
        )
        results.append(result)
        print_test_result(result)
        
        # Test 22: Add Favorite Verse (Duplicate)
        result = test_endpoint(
            "Add Favorite Verse - Duplicate",
            "POST",
            "/api/users/me/favorites",
            headers=auth_headers,
            body={
                "book_id": 1,
                "chapter": 1,
                "verse": 1
            }
        )
        results.append(result)
        print_test_result(result)
        
        # Test 23: Add Favorite Verse (Invalid Verse)
        result = test_endpoint(
            "Add Favorite Verse - Invalid Verse",
            "POST",
            "/api/users/me/favorites",
            headers=auth_headers,
            body={
                "book_id": 999,
                "chapter": 999,
                "verse": 999
            }
        )
        results.append(result)
        print_test_result(result)
        
        # Test 24: Get Favorite Verses
        result = test_endpoint(
            "Get Favorite Verses",
            "GET",
            "/api/users/me/favorites?page=1&limit=10",
            headers={"Authorization": f"Bearer {auth_token}"}
        )
        results.append(result)
        print_test_result(result)
        
        # Test 25: Get Favorite Verses (Pagination)
        result = test_endpoint(
            "Get Favorite Verses - Pagination",
            "GET",
            "/api/users/me/favorites?page=1&limit=5",
            headers={"Authorization": f"Bearer {auth_token}"}
        )
        results.append(result)
        print_test_result(result)
        
        # Test 26: Add Highlighted Verse
        result = test_endpoint(
            "Add Highlighted Verse",
            "POST",
            "/api/users/me/highlights",
            headers=auth_headers,
            body={
                "book_id": 1,
                "chapter": 1,
                "verse": 2,
                "note": "This is a test highlight",
                "color": "yellow"
            }
        )
        results.append(result)
        print_test_result(result)
        
        # Test 27: Add Highlighted Verse (Duplicate)
        result = test_endpoint(
            "Add Highlighted Verse - Duplicate",
            "POST",
            "/api/users/me/highlights",
            headers=auth_headers,
            body={
                "book_id": 1,
                "chapter": 1,
                "verse": 2,
                "note": "Updated note",
                "color": "blue"
            }
        )
        results.append(result)
        print_test_result(result)
        
        # Test 28: Get Highlighted Verses
        result = test_endpoint(
            "Get Highlighted Verses",
            "GET",
            "/api/users/me/highlights?page=1&limit=10",
            headers={"Authorization": f"Bearer {auth_token}"}
        )
        results.append(result)
        print_test_result(result)
        
        # Test 29: Update Highlighted Verse
        result = test_endpoint(
            "Update Highlighted Verse",
            "PUT",
            "/api/users/me/highlights/1/1/2",
            headers=auth_headers,
            body={
                "note": "Updated highlight note",
                "color": "green"
            }
        )
        results.append(result)
        print_test_result(result)
        
        # Test 30: Update Last Read Position
        result = test_endpoint(
            "Update Last Read Position",
            "POST",
            "/api/users/me/last-read",
            headers=auth_headers,
            body={
                "book_id": 1,
                "book_name": "Genesis",
                "chapter": 1,
                "verse": 5
            }
        )
        results.append(result)
        print_test_result(result)
        
        # Test 31: Get Last Read Position
        result = test_endpoint(
            "Get Last Read Position",
            "GET",
            "/api/users/me/last-read",
            headers={"Authorization": f"Bearer {auth_token}"}
        )
        results.append(result)
        print_test_result(result)
        
        # Test 32: Get Last Read Verses (Legacy Endpoint)
        result = test_endpoint(
            "Get Last Read Verses (Legacy)",
            "GET",
            "/api/last-read-verses/",
            headers={"Authorization": f"Bearer {auth_token}"}
        )
        results.append(result)
        print_test_result(result)
        
        # Test 33: Remove Favorite Verse
        result = test_endpoint(
            "Remove Favorite Verse",
            "DELETE",
            "/api/users/me/favorites/1/1/1",
            headers={"Authorization": f"Bearer {auth_token}"}
        )
        results.append(result)
        print_test_result(result)
        
        # Test 34: Remove Highlighted Verse
        result = test_endpoint(
            "Remove Highlighted Verse",
            "DELETE",
            "/api/users/me/highlights/1/1/2",
            headers={"Authorization": f"Bearer {auth_token}"}
        )
        results.append(result)
        print_test_result(result)
        
        # Test 35: Get Current User Profile (After Update)
        result = test_endpoint(
            "Get Current User Profile - After Update",
            "GET",
            "/api/users/me",
            headers={"Authorization": f"Bearer {auth_token}"}
        )
        results.append(result)
        print_test_result(result)
        
    else:
        print(f"\n{Colors.YELLOW}⚠ Warning: No auth token available, skipping protected endpoint tests{Colors.RESET}")
    
    # Test 36: Access Protected Endpoint Without Token
    result = test_endpoint(
        "Get Current User Profile - No Token",
        "GET",
        "/api/users/me"
    )
    results.append(result)
    print_test_result(result)
    
    # Test 37: Access Protected Endpoint With Invalid Token
    result = test_endpoint(
        "Get Current User Profile - Invalid Token",
        "GET",
        "/api/users/me",
        headers={"Authorization": "Bearer invalid_token_12345"}
    )
    results.append(result)
    print_test_result(result)
    
    # Test 38: Update Last Read (Invalid Verse)
    if auth_token:
        result = test_endpoint(
            "Update Last Read - Invalid Verse",
            "POST",
            "/api/users/me/last-read",
            headers={
                "Authorization": f"Bearer {auth_token}",
                "Content-Type": "application/json"
            },
            body={
                "book_id": 999,
                "book_name": "InvalidBook",
                "chapter": 999,
                "verse": 999
            }
        )
        results.append(result)
        print_test_result(result)
    
    # Generate and save report
    print(f"\n{Colors.BOLD}{Colors.BLUE}")
    print("=" * 60)
    print("Generating Report...")
    print("=" * 60)
    print(f"{Colors.RESET}")
    
    report = generate_report(results)
    report_filename = f"API_TEST_REPORT_{datetime.now().strftime('%Y%m%d_%H%M%S')}.md"
    report_file = os.path.join(REPORTS_DIR, report_filename)
    
    with open(report_file, 'w') as f:
        f.write(report)
    
    print(f"\n{Colors.GREEN}✓ Report saved to: {report_file}{Colors.RESET}")
    
    # Print summary
    total = len(results)
    passed = sum(1 for r in results if r.success)
    failed = total - passed
    
    print(f"\n{Colors.BOLD}Test Summary:{Colors.RESET}")
    print(f"  Total: {total}")
    print(f"  {Colors.GREEN}Passed: {passed}{Colors.RESET}")
    print(f"  {Colors.RED}Failed: {failed}{Colors.RESET}")
    print(f"  Success Rate: {(passed/total*100):.1f}%")
    
    if failed > 0:
        print(f"\n{Colors.YELLOW}Note: Some tests failed. Check the report for details.{Colors.RESET}")
        return 1
    
    return 0

if __name__ == "__main__":
    sys.exit(main())


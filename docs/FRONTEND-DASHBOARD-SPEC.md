# HIAS Core — Frontend Admin Dashboard Specification

## Overview

Build a Next.js admin dashboard for the HIAS Core health insurance administration system. The backend is a Go REST API with 246+ endpoints covering 12 domains: Identity, Product, Policy, Claims, Billing, Provider, Sales, PreAuth, Notification, Audit, Analytics, and Reinsurance.

This document is a **complete specification** — it contains everything needed to build the frontend without reading the backend source code.

---

## 1. Project Setup

### Tech Stack

| Technology | Purpose |
|-----------|---------|
| Next.js 14+ (App Router) | React framework with SSR |
| TypeScript | Type safety |
| Tailwind CSS | Utility-first styling |
| shadcn/ui | Component library (built on Radix UI) |
| TanStack Query (React Query) | Server state management |
| Axios | HTTP client |
| Zustand | Client state (auth, UI preferences) |
| React Hook Form + Zod | Form handling + validation |
| date-fns | Date formatting/manipulation |
| Recharts | Charts and data visualization |
| Lucide React | Icons |
| nuqs | URL query state management |

### Project Structure

```
src/
├── app/
│   ├── (auth)/
│   │   ├── login/page.tsx
│   │   └── register/page.tsx
│   ├── (dashboard)/
│   │   ├── layout.tsx                  # Dashboard layout (sidebar + topbar)
│   │   ├── page.tsx                    # Dashboard home
│   │   ├── products/
│   │   │   ├── plans/page.tsx
│   │   │   ├── plans/[id]/page.tsx
│   │   │   ├── plans/new/page.tsx
│   │   │   └── providers/page.tsx
│   │   ├── sales/
│   │   │   ├── leads/page.tsx
│   │   │   ├── leads/[id]/page.tsx
│   │   │   ├── quotations/page.tsx
│   │   │   ├── quotations/[id]/page.tsx
│   │   │   ├── quotations/new/page.tsx
│   │   │   └── approval-limits/page.tsx
│   │   ├── policies/
│   │   │   ├── page.tsx
│   │   │   ├── [id]/page.tsx
│   │   │   ├── new/page.tsx
│   │   │   ├── endorsements/page.tsx
│   │   │   ├── renewals/page.tsx
│   │   │   ├── underwriting/page.tsx
│   │   │   └── credit-notes/page.tsx
│   │   ├── claims/
│   │   │   ├── page.tsx
│   │   │   ├── [id]/page.tsx
│   │   │   ├── new/page.tsx
│   │   │   ├── preauths/page.tsx
│   │   │   ├── preauths/[id]/page.tsx
│   │   │   ├── cases/page.tsx
│   │   │   └── cases/[id]/page.tsx
│   │   ├── billing/
│   │   │   ├── invoices/page.tsx
│   │   │   ├── payments/page.tsx
│   │   │   └── remittances/page.tsx
│   │   ├── reinsurance/
│   │   │   ├── page.tsx               # Reinsurance dashboard
│   │   │   ├── treaties/page.tsx
│   │   │   ├── treaties/[id]/page.tsx
│   │   │   ├── treaties/new/page.tsx
│   │   │   ├── cessions/page.tsx
│   │   │   ├── recoveries/page.tsx
│   │   │   ├── bordereaux/page.tsx
│   │   │   └── statements/page.tsx
│   │   ├── providers/
│   │   │   ├── page.tsx
│   │   │   ├── [id]/page.tsx
│   │   │   └── new/page.tsx
│   │   ├── users/
│   │   │   ├── page.tsx
│   │   │   └── [id]/page.tsx
│   │   ├── notifications/page.tsx
│   │   ├── audit/page.tsx
│   │   └── analytics/page.tsx
│   └── layout.tsx                      # Root layout
├── components/
│   ├── ui/                            # shadcn/ui components
│   ├── layout/
│   │   ├── sidebar.tsx
│   │   ├── topbar.tsx
│   │   └── breadcrumbs.tsx
│   ├── data-table/
│   │   ├── data-table.tsx             # Generic server-side paginated table
│   │   ├── column-header.tsx
│   │   ├── pagination.tsx
│   │   └── toolbar.tsx
│   ├── forms/
│   │   ├── plan-form.tsx
│   │   ├── claim-form.tsx
│   │   ├── member-form.tsx
│   │   ├── quotation-form.tsx
│   │   └── ...
│   ├── status-badge.tsx
│   ├── money-display.tsx
│   ├── confirm-dialog.tsx
│   ├── file-upload.tsx
│   └── charts/
│       ├── bar-chart.tsx
│       ├── line-chart.tsx
│       └── pie-chart.tsx
├── lib/
│   ├── api/
│   │   ├── client.ts                  # Axios instance with interceptors
│   │   ├── auth.ts                    # Auth endpoints
│   │   ├── users.ts
│   │   ├── plans.ts
│   │   ├── benefits.ts
│   │   ├── claims.ts
│   │   ├── policies.ts
│   │   ├── members.ts
│   │   ├── quotations.ts
│   │   ├── leads.ts
│   │   ├── providers.ts
│   │   ├── billing.ts
│   │   ├── reinsurance.ts
│   │   ├── notifications.ts
│   │   ├── audit.ts
│   │   └── analytics.ts
│   ├── hooks/
│   │   ├── use-auth.ts
│   │   ├── use-claims.ts
│   │   ├── use-policies.ts
│   │   └── ...                        # TanStack Query hooks per domain
│   ├── types/
│   │   ├── auth.ts
│   │   ├── product.ts
│   │   ├── policy.ts
│   │   ├── claims.ts
│   │   ├── billing.ts
│   │   ├── sales.ts
│   │   ├── provider.ts
│   │   ├── reinsurance.ts
│   │   ├── notification.ts
│   │   ├── audit.ts
│   │   ├── analytics.ts
│   │   └── common.ts
│   ├── utils/
│   │   ├── money.ts                   # formatMoney(800000) → "KES 8,000.00"
│   │   ├── date.ts                    # Date formatting helpers
│   │   ├── status-colors.ts           # Status → color mapping
│   │   └── constants.ts               # Enum values, roles
│   └── stores/
│       ├── auth-store.ts              # Zustand auth store
│       └── ui-store.ts                # Sidebar collapsed, theme
├── middleware.ts                       # Route protection
└── providers/
    ├── query-provider.tsx
    └── auth-provider.tsx
```

---

## 2. Authentication Module

### Token Management

```typescript
// lib/stores/auth-store.ts
interface AuthState {
  accessToken: string | null;
  refreshToken: string | null;
  user: UserResponse | null;
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  refreshAccessToken: () => Promise<void>;
}
```

**Implementation requirements:**
- Store tokens in `httpOnly` cookies or `localStorage` (cookies preferred for SSR)
- Access token included in every API request via Axios interceptor: `Authorization: Bearer <token>`
- On 401 response, attempt token refresh via `/api/v1/auth/refresh`
- If refresh fails, redirect to `/login`
- On login success, store user object and redirect to dashboard

### Login Page (`/login`)
- Email + password form
- "Forgot password" link (placeholder)
- Redirect to dashboard on success
- Show error toast on invalid credentials

### Register Page (`/register`)
- Name, email, phone, password fields
- Redirect to login on success

### Route Protection (`middleware.ts`)
- Check for valid token on all `/dashboard/*` routes
- Redirect to `/login` if unauthenticated
- Role-based page access (see Section 8)

---

## 3. Layout & Navigation

### Dashboard Layout

```
┌─────────────────────────────────────────────────────────┐
│  HIAS Logo   │  Breadcrumbs                 🔔 3  👤 Admin ▾  │
├──────────────┼──────────────────────────────────────────┤
│              │                                          │
│  Dashboard   │     Page Content                        │
│              │                                          │
│  Products ▾  │                                          │
│   Plans      │                                          │
│   Providers  │                                          │
│              │                                          │
│  Sales ▾     │                                          │
│   Leads      │                                          │
│   Quotations │                                          │
│   Approvals  │                                          │
│              │                                          │
│  Policies ▾  │                                          │
│   All        │                                          │
│   Endorsmnts │                                          │
│   Renewals   │                                          │
│   UW Flags   │                                          │
│   Credit Nts │                                          │
│              │                                          │
│  Claims ▾    │                                          │
│   All        │                                          │
│   Pre-Auths  │                                          │
│   Cases      │                                          │
│              │                                          │
│  Billing ▾   │                                          │
│   Invoices   │                                          │
│   Payments   │                                          │
│   Remittances│                                          │
│              │                                          │
│  Reinsurance▾│                                          │
│   Dashboard  │                                          │
│   Treaties   │                                          │
│   Cessions   │                                          │
│   Recoveries │                                          │
│   Bordereaux │                                          │
│   Statements │                                          │
│              │                                          │
│  Providers   │                                          │
│  Users       │                                          │
│  Audit Trail │                                          │
│  Analytics   │                                          │
│              │                                          │
└──────────────┴──────────────────────────────────────────┘
```

### Top Bar
- Left: Breadcrumbs (auto-generated from route)
- Right: Notification bell with unread count badge (poll `GET /notifications/unread-count` every 30s)
- Right: User avatar/name dropdown: Profile, Logout

### Sidebar
- Collapsible (icon-only mode)
- Grouped sections with expand/collapse
- Active page highlighted
- Role-based visibility (see Section 8)
- Mobile: overlay drawer

---

## 4. Dashboard Pages by Scenario

### 4.1 Main Dashboard (`/dashboard`)

**KPI Cards Row:**
| Card | API Source | Display |
|------|----------|---------|
| Total Policies | Custom count | Number |
| Active Claims | Claims volume | Number |
| Loss Ratio | `GET /analytics/kpis` | Percentage with trend arrow |
| Premium Collected | `GET /analytics/kpis` | Money (KES) |
| Approval Rate | `GET /analytics/kpis` | Percentage |
| Fraud Rate | `GET /analytics/kpis` | Percentage (red if > 5%) |

**Charts Row:**
- Bar chart: Claims volume by month (approved vs rejected vs pending)
- Pie chart: Loss ratio breakdown by benefit category
- Line chart: Premium trend (monthly)

**Quick Actions:**
- "New Claim" button
- "New Policy" button
- "New Quotation" button

**Recent Activity Feed:**
- Latest 10 audit events (`GET /audit`)
- Display: timestamp, user, action, entity

---

### 4.2 Product Configuration

#### Plans List (`/products/plans`)
- Data table: Name, Type, Segment, Base Premium, Status, Created
- Filters: type (individual/group), segment, status
- "New Plan" button
- Click row → plan detail page

#### Plan Detail (`/products/plans/[id]`)
- **Header**: Plan name, type, segment, base premium, status toggle
- **Tabs:**
  - **Benefits**: Tree view showing parent benefits and sub-benefits
    - Each benefit row: name, category, annual limit, co-pay, waiting period, deductible
    - "Add Benefit" button → modal form
    - Expand to see sub-benefits; "Add Sub-Benefit" button per parent
  - **Exclusions**: Table of exclusions (description, type, ICD codes)
    - "Add Exclusion" button → modal form
    - Inline edit/delete
  - **Premium Rules**: Rate sheet table
    - Columns: calculation_type, relationship, rate_amount, age range, discount
    - "Add Rule" button → modal form
    - "Calculate Premium" button → modal with relationship/DOB inputs → shows calculated premium
  - **Underwriting Rules**: Table of rules
    - Columns: rule_type, relationship, parameter, severity, blocking, active
    - "Add Rule" button → modal form
    - Toggle active/inactive
  - **Provider Network**: Table of associated providers
    - Columns: provider name, benefit category, status
    - "Add Provider" button → provider search/select modal
    - Toggle status (ACTIVE/INACTIVE)

---

### 4.3 Business Development

#### Leads Pipeline (`/sales/leads`)
- **Default view**: Kanban board
  - Columns: NEW → CONTACTED → QUALIFIED → PROPOSAL_SENT → NEGOTIATION → WON / LOST
  - Cards show: contact name, company, expected premium, probability, next follow-up
  - Drag-and-drop to change status (calls `PUT /leads/:id/status`)
- **Alternative view**: Data table (toggle button)
  - Columns: Lead#, Contact, Company, Source, Expected Premium, Status, Follow-up Date
  - Filter by status, source, segment
- "New Lead" button → form modal

#### Lead Detail (`/sales/leads/[id]`)
- **Header**: Contact info, company, source, segment, expected premium, status badge
- **Sections:**
  - **Activity Timeline**: Chronological feed of activities (calls, emails, meetings, notes)
    - "Add Activity" button → modal with type, description, scheduled/completed dates
  - **Quotations**: Table of quotations linked to this lead
    - "Create Quotation" button → redirect to new quotation page with lead pre-filled
  - **Info Panel**: Assigned agent, closure probability, next follow-up date (editable)

#### Quotations List (`/sales/quotations`)
- Data table: Quotation#, Client, Plan, Type, Status, Premium, Version, Created
- Filters: status, quotation_type
- "New Quotation" button

#### New Quotation Wizard (`/sales/quotations/new`)
- **Step 1**: Select lead (search/dropdown) + Select plan
- **Step 2**: Enter member count + proposed members (relationship + DOB per member)
- **Step 3**: Select billing frequency. System shows calculated base premium.
- **Step 4**: Optional discount (type + value + reason) + optional loading (type + value + reason)
- **Step 5**: Review summary showing:
  - Base premium, discount amount, loading amount, final premium
  - If exceeds approval limits → shows "Requires Approval" warning
- **Submit**: Creates quotation + version 1

#### Quotation Detail (`/sales/quotations/[id]`)
- **Header**: Quotation#, client name, plan, status badge, action buttons
- **Action Buttons** (based on status):
  - DRAFT: "Issue Quotation"
  - ISSUED: "Send to Client" (channel select: SMS/Email)
  - PENDING_DECISION: "Accept" / "Decline"
  - ACCEPTED: "Convert to Policy" (date picker for start date)
- **Tabs:**
  - **Versions**: Table of all versions
    - Columns: Version#, Premium, Discount, Loading, Final Premium, Approval Status
    - "New Version" button → same form as wizard steps 2-5
    - "Compare" button → side-by-side comparison (select two versions)
      - Shows: premium diff, discount diff, loading diff, member count diff
    - Per version: "Submit for Approval" / "Approve" / "Reject" buttons (role-based)
  - **Documents**: File list with upload/download/delete
    - Drag-and-drop upload area
    - Table: filename, type, size, uploaded by, date

#### Approval Limits (`/sales/approval-limits`) — Admin only
- Data table: Role, Max Discount %, Max Discount Amount, Max Loading %, Max Loading Amount, Escalation Role
- "Add Limit" button
- Inline edit

---

### 4.4 Policy Lifecycle

#### Policies List (`/policies`)
- Data table: Policy#, Policyholder, Plan, Status, Premium, Start Date, End Date
- **Status tabs**: ALL | DRAFT | ACTIVE | LAPSED | SUSPENDED | TERMINATED
  - Uses `GET /policies/by-status?status=ACTIVE` for tab content
- Filters: date range, plan
- "New Policy" button
- Bulk actions: "Activate Selected" / "Lapse Selected" (checkboxes)

#### Policy Detail (`/policies/[id]`)
- **Header**: Policy#, policyholder name/email/phone, status badge, premium amount, dates
- **Action Buttons** (based on status):
  - DRAFT: "Activate" (requires payment reference input)
  - ACTIVE: "Suspend" / "Lapse" / "Terminate" / "Change Plan"
  - LAPSED/SUSPENDED: "Reinstate"
- **Tabs:**
  - **Members**:
    - Table: Member#, Name, Relationship, DOB, Gender, Status, Verified
    - "Enroll Member" button → form modal
    - "Bulk Enroll" button → multi-member form
    - "Import CSV" button → file upload
    - "Bulk Remove" button → checkbox select + confirm dialog
    - Click member → sub-panel with:
      - Member details, eligibility status, underwriting flags
      - "Verify", "Suspend", "Reactivate", "Remove" buttons
      - "Generate Card" button
  - **Endorsements**:
    - Table: Type, Status, Effective Date, Premium Adjustment, Requested By
    - "Create Endorsement" button → form (type select, effective date, changes JSON, reason)
    - Per endorsement: "Approve" / "Reject" / "Apply" buttons
  - **Documents**:
    - Table: Type, Filename, Size, Generated By, Date
    - Action buttons: "Generate Welcome Letter", "Generate Policy Schedule", "Generate Member Cards"
  - **Billing**:
    - **Installments sub-tab**: Schedule table + individual installments with "Mark Paid" buttons
    - **Invoices sub-tab**: List from policy invoices
    - **Credit Notes sub-tab**: List with "Approve" / "Apply" buttons
  - **Renewals**:
    - Table: Status, Renewal Date, New Premium, Plan Change, Approved By
    - "Initiate Renewal" button → form (renewal date, optional new plan)
    - Per renewal: "Approve" / "Reject" / "Complete" buttons
  - **Underwriting**:
    - Assessments table: Status, Risk Score, Decision
    - "Submit Assessment" button → questionnaire form
    - Flags table: Type, Severity, Status, Details
    - Per flag: "Resolve" / "Override" buttons (Underwriter/Admin)
  - **Cases**:
    - Table from `GET /policies/:id/cases`

#### Underwriting Flags Dashboard (`/policies/underwriting`)
- Table: Flag#, Policy#, Member, Flag Type, Severity, Status, Created
- Filter by: status (OPEN, RESOLVED, OVERRIDDEN), severity
- Badge showing open count (`GET /underwriting-flags/count`)
- Click → flag detail with resolve/override actions

---

### 4.5 Claims Processing

#### Claims Dashboard (`/claims`)
- **KPI Cards** (from `GET /analytics/dashboard`):
  - Total Claims, Approved, Rejected, Manual Review, Paid
  - Average TAT, Approval Rate, SLA Breach Count
- **Status Tabs**: ALL | RECEIVED | VALIDATED | APPROVED | VETTED | READY_FOR_PAYMENT | PAID | REJECTED | MANUAL_REVIEW
- Data table: Claim#, Policy#, Member, Provider, Amount, Status, Service Date, SLA
  - **SLA column**: Shows time remaining or "BREACHED" in red
  - Rows with breached SLA highlighted in red background
- "New Claim" button
- "Import CSV" button (ClaimsOfficer/Admin)
- "Bulk Submit" button

#### Claim Detail (`/claims/[id]`)
- **Header**: Claim#, status badge, total amount, approved amount, claim type, service date
- **SLA indicator**: Countdown or breach warning
- **Action Buttons** (role-based, status-dependent):
  - APPROVED/MANUAL_REVIEW: "Vet" (ClaimsOfficer) → vetted amount input
  - ADJUDICATED/MANUAL_REVIEW: "Approve" / "Reject" (Manager)
  - VETTED/PARTIALLY_VETTED: "Ready for Payment" (Finance)
  - READY_FOR_PAYMENT: "Mark Paid" / "Mark Part Paid" (Finance)
  - REJECTED: "Generate Decline Letter"
- **Sections:**
  - **Line Items Table**: Procedure code, name, quantity, unit price, total, approved amount
  - **Adjudication Decision Panel** (if exists):
    - Decision: APPROVE/REJECT/MANUAL_REVIEW
    - Payable amount, member responsibility, deductible applied, co-pay applied
    - Benefit category
    - Rule results: expandable list showing each rule, category, status (PASS/FAIL/FLAG), reason
  - **Fraud Flags** (if any):
    - Table: Flag Type, Severity, Details, Resolved
    - Color-coded by severity: LOW=yellow, MEDIUM=orange, HIGH=red, CRITICAL=purple
  - **Documents**:
    - Upload area + file list
  - **Audit Trail**: Recent events for this claim

#### New Claim Form (`/claims/new`)
- **Fields:**
  - Policy (search by policy#) → auto-loads members
  - Member (select from policy members)
  - Provider (search by name)
  - Pre-Auth (optional, search by auth code)
  - Claim Type: DIRECT | REIMBURSEMENT | CREDIT | EXCEPTION
  - Diagnosis Codes (multi-input tags)
  - Service Date, Admission Date (optional), Discharge Date (optional)
  - Notes
- **Line Items** (dynamic table):
  - Add row: Procedure Code, Procedure Name, Diagnosis Code, Quantity, Unit Price
  - Auto-calculate: Total Price = Quantity × Unit Price
  - Running total displayed
- **Submit** → runs full pipeline, shows result with adjudication decision

#### Pre-Authorizations (`/claims/preauths`)
- Data table: Auth Code, Policy#, Member, Provider, Estimated Cost, Approved Amount, Status
- "New Pre-Auth" button
- Filters: status

#### Pre-Auth Detail (`/claims/preauths/[id]`)
- Header: Auth code, status, estimated cost, approved amount, validity period
- Action buttons:
  - SUBMITTED: "Review" → decision form (APPROVED/DENIED/INFO_REQUESTED)
  - APPROVED: "Generate LOU"
- Procedure codes and diagnosis codes display
- Related cases list

#### Cases (`/claims/cases`)
- Data table: Case#, Policy#, Member, Provider, Status, Admission Date, Diagnosis
- Filter by status: SCHEDULED | ADMITTED | IN_TREATMENT | DISCHARGED | CLOSED
- Status count badges (`GET /cases/count`)

#### Case Detail (`/claims/cases/[id]`)
- Header: Case#, status, estimated cost, actual cost
- Action buttons (status-dependent):
  - SCHEDULED: "Admit" (admission date input)
  - ADMITTED: "Start Treatment"
  - IN_TREATMENT: "Discharge" (discharge date + actual cost inputs)
  - DISCHARGED: "Close Case"
- Info: Pre-auth, treating doctor, room type, diagnosis
- Timeline of status transitions

---

### 4.6 Billing

#### Invoices (`/billing/invoices`)
- Data table: Invoice#, Policy#, Amount, Due Date, Status, Period
- Status filter: PENDING | PAID | OVERDUE | CANCELLED
- Overdue invoices highlighted in red

#### Payments (`/billing/payments`)
- Data table: Reference#, Type, Amount, Method, Status, Date
- "New Payment" button → form (invoice select, amount, method, phone for MPESA)
- Per payment: "Retry" button (for FAILED), "Reconcile" button (for CONFIRMED)

#### Remittances (`/billing/remittances`)
- Data table: Provider, Amount, Status, Period, Advice Sent
- "Create Remittance" button → form (provider select, period)
- Per remittance: "Export Payment File" button (downloads CSV/PDF)

---

### 4.7 Reinsurance

#### Reinsurance Dashboard (`/reinsurance`)
- **KPI Cards** (from `GET /analytics/reinsurance`):
  - Active Treaties, Total Ceded Premiums, Total Recoverable, Total Recovered
  - Outstanding, Cession Ratio (%), Recovery Success Rate (%)
  - Unacknowledged Alerts (with link to alerts)

#### Treaties List (`/reinsurance/treaties`)
- Data table: Treaty#, Name, Type, Status, Effective Date, Expiry Date, Retention Limit
- Filter by: type (QUOTA_SHARE/XOL), status
- "New Treaty" button

#### Treaty Detail (`/reinsurance/treaties/[id]`)
- **Header**: Treaty#, name, type badge, status badge, dates, retention limit
- **Action Buttons**:
  - DRAFT: "Activate" (validates participants/shares)
  - ACTIVE: "Terminate"
- **Tabs:**
  - **Participants**:
    - Table: Reinsurer Name, Share %, Commission Rate, Lead Reinsurer
    - Total share percentage displayed (must be ≤ 100%)
    - "Add Participant" button → form modal
    - Inline edit/delete
  - **Layers** (for XOL treaties):
    - Table: Layer#, Attachment Point, Layer Limit, Deductible, Premium Rate, Aggregate Limit, Aggregate Used
    - "Add Layer" button → form modal
    - Progress bar showing aggregate usage (green/yellow/red)
  - **Profit Commission Rules**:
    - Table: Type, Loss Ratio Range, Commission Rate, Carry Forward Years
    - "Add Rule" button → form modal
    - "Calculate Profit Commission" button → period selector → shows result
  - **Cessions**: Cessions filtered to this treaty
  - **Recoveries**: Recoveries filtered to this treaty
  - **Bordereaux**: Bordereaux filtered to this treaty
  - **Statements**: Statements filtered to this treaty
  - **Alerts**: Alerts filtered to this treaty
    - Per alert: "Acknowledge" button
    - Color-coded by severity

#### Cessions (`/reinsurance/cessions`)
- Data table: Cession#, Treaty, Policy, Gross, Ceded, Retained, Commission, Status
- "Cede Premium" button → form (treaty select, policy select, amount)
- "Auto-Cede" button → form (policy select, amount) — cedes to all quota share treaties
- Per cession: "Book" / "Reverse" buttons

#### Recoveries (`/reinsurance/recoveries`)
- Data table: Recovery#, Claim#, Treaty, Gross, Recoverable, Recovered, Outstanding, Status
- "Create Recovery" button → form
- "Apply for Claim" button → claim select, approved amount
- Filter tabs: ALL | OUTSTANDING | PAID | WRITTEN_OFF
- **Aged Analysis** view: Bar chart by aging bucket (0-30, 31-60, 61-90, 90+ days)
- Per recovery: Workflow actions (Acknowledge → Request Info → Approve → Record Payment / Write Off)

#### Recovery Detail (modal or page):
- Header: Recovery#, status, amounts
- Workflow timeline: Visual stepper showing NOTIFIED → ACKNOWLEDGED → APPROVED → PAID
- Workflow event log: Table of all transitions with notes and timestamps
- Action buttons based on current status

#### Bordereaux (`/reinsurance/bordereaux`)
- Data table: Bordereau#, Treaty, Type, Period, Total Gross, Total Ceded, Status
- "Generate Premium Bordereau" button → form (treaty, period)
- "Generate Claim Bordereau" button → form (treaty, period)
- Per bordereau: "Finalize" / "Mark Sent" buttons
- Click → item detail table

#### Reinsurer Statements (`/reinsurance/statements`)
- Data table: Statement#, Treaty, Participant, Premium Ceded, Claims Recovered, Commission, Net Balance, Status
- "Generate Statement" button → form (treaty, participant, period)
- "Calculate Profit Commission" button → form (treaty, period) → shows calculation result
- Per statement: "Issue" / "Acknowledge" / "Settle" buttons

#### Treaty Alerts (shown in treaty detail + separate panel)
- Alert list: Type, Severity, Message, Threshold vs Current, Acknowledged
- "Check Limits" button → triggers limit check for a treaty
- "Check Catastrophe" button → triggers catastrophe check
- "Check Expiry" button → checks all treaties
- Unacknowledged count badge in sidebar
- Per alert: "Acknowledge" button

---

### 4.8 Provider Management

#### Providers List (`/providers`)
- Data table: Name, Type, License#, Status, Tier, Accreditation, County
- Filters: type, status, tier, accreditation status
- "Register Provider" button

#### Provider Detail (`/providers/[id]`)
- **Header**: Name, type, license, status badge, tier badge, accreditation status
- **Action Buttons** (based on status):
  - PENDING: "Start Credentialing"
  - CREDENTIALING: "Activate"
  - ACTIVE: "Suspend" / "Terminate"
  - SUSPENDED: "Activate" / "Terminate"
- **Tabs:**
  - **Details**: Editable form (name, contact, address, email, phone)
  - **Tier**: Current tier + "Update Tier" dropdown
  - **Accreditation**: Status, expiry date, body + "Update Accreditation" form
  - **Contracts**: Table (start/end dates, terms, status) + "Add Contract" button
  - **Rate Cards**: Table (procedure code, name, rate, age range) + "Add Rate Card" / "Bulk Import" buttons
  - **Statements**: Provider statements with reconciliation
  - **Cases**: Inpatient cases at this provider

---

### 4.9 User Management (Admin Only)

#### Users List (`/users`)
- Data table: Name, Email, Phone, Role, Status
- "Create User" button → form modal
- Per user: "Edit" / "Assign Role" / "Update Status" actions

---

### 4.10 Cross-Cutting Features

#### Notifications (`/notifications`)
- List view: Subject, body, channel, type, status (read/unread), date
- "Mark as Read" button per notification
- Unread highlighted with dot
- Bell icon in topbar with unread count badge

#### Audit Trail (`/audit`)
- Data table: Timestamp, User, Entity Type, Entity ID, Action, Details
- Filters: entity type, action, date range, user
- Expandable rows showing old_value / new_value JSON diff

#### Analytics (`/analytics`)
- **Claims Volume Chart**: Bar chart (approved/rejected/pending by month)
- **Premium Trend Chart**: Line chart (monthly premium collected)
- **Loss Ratio Chart**: Pie chart (breakdown by category)
- **KPI Cards**: Approval rate, average TAT, loss ratio, fraud rate
- **Top Providers Table**: Name, claim count, total amount, total approved
- "Export CSV" button

---

## 5. TypeScript Types

### Common Types

```typescript
// lib/types/common.ts

interface ApiResponse<T> {
  status: "success" | "error";
  message: string;
  data: T;
}

interface PaginatedResponse<T> {
  status: "success";
  message: string;
  data: T[];
  page: number;
  page_size: number;
  total_count: number;
  total_pages: number;
}

interface ErrorResponse {
  status: "error";
  message: string;
  error?: string;
}

interface BulkResult {
  succeeded: number;
  failed: number;
  errors?: string[];
}

type UUID = string;
type Money = number; // cents — display with formatMoney()
type ISODateTime = string; // "2026-01-15T10:30:00Z"
```

### All Enum / Status Types (Complete Reference)

Every backend enum mapped to TypeScript. These MUST be used in interfaces and UI components (status badges, select dropdowns, filter tabs).

```typescript
// lib/types/enums.ts — COMPLETE BACKEND ENUM MAPPING (75 types)

// ── Identity & Auth ─────────────────────────────────────────
export type UserRole = "Admin" | "Underwriter" | "ClaimsOfficer" | "Finance" | "Provider" | "Member" | "SalesAgent" | "Manager";
export type UserStatus = "ACTIVE" | "INACTIVE" | "SUSPENDED";

// ── Policy Management ───────────────────────────────────────
export type PolicyStatus = "DRAFT" | "ACTIVE" | "LAPSED" | "TERMINATED" | "SUSPENDED";
export type MemberStatus = "ACTIVE" | "SUSPENDED" | "REMOVED";
export type MemberRelationship = "principal" | "spouse" | "child" | "parent";
export type Gender = "male" | "female" | "other";
export type EndorsementType = "ADD_MEMBER" | "REMOVE_MEMBER" | "UPDATE_MEMBER" | "PLAN_CHANGE";
export type EndorsementStatus = "PENDING" | "APPROVED" | "REJECTED" | "APPLIED";
export type RenewalStatus = "PENDING" | "APPROVED" | "REJECTED" | "COMPLETED" | "EXPIRED";
export type PolicyDocumentType = "WELCOME_LETTER" | "MEMBER_CARD" | "POLICY_SCHEDULE" | "RENEWAL_NOTICE" | "ENDORSEMENT" | "LOU" | "DECLINE_LETTER";
export type CreditNoteStatus = "DRAFT" | "APPROVED" | "APPLIED" | "CANCELLED";

// ── Underwriting ────────────────────────────────────────────
export type UnderwritingStatus = "PENDING" | "APPROVED" | "DECLINED" | "REFER";
export type UnderwritingRuleType = "MAX_AGE" | "MIN_AGE" | "DOUBLE_INSURANCE" | "PRE_EXISTING_CONDITION" | "BMI_THRESHOLD" | "WAITING_PERIOD";
export type UnderwritingFlagType = "MAX_AGE" | "MIN_AGE" | "DOUBLE_INSURANCE" | "PRE_EXISTING_CONDITION" | "BMI_THRESHOLD" | "WAITING_PERIOD" | "RENEWAL_SKIP";
export type UnderwritingFlagStatus = "OPEN" | "ACKNOWLEDGED" | "RESOLVED" | "OVERRIDDEN";
export type UnderwritingSeverity = "LOW" | "MEDIUM" | "HIGH";

// ── Product Configuration ───────────────────────────────────
export type PlanType = "individual" | "group";
export type PlanSegment = "retail" | "corporate" | "sme";
export type PlanStatus = "ACTIVE" | "INACTIVE";
export type BenefitCategory = "outpatient" | "inpatient" | "dental" | "optical" | "maternity";
export type CoPayType = "percentage" | "fixed";
export type SubLimitType = "none" | "per_visit" | "per_item";
export type WaitingPeriodType = "general" | "maternity" | "pre_existing" | "chronic" | "surgical";
export type ExclusionType = "pre_existing" | "cosmetic" | "experimental";
export type PremiumCalculationType = "flat" | "per_member" | "tiered" | "per_family";
export type DiscountType = "percentage" | "fixed";
export type LoadingType = "percentage" | "fixed";
export type BillingFrequency = "monthly" | "quarterly" | "semi_annual" | "annual";

// ── Claims Processing ───────────────────────────────────────
export type ClaimStatus = "RECEIVED" | "VALIDATED" | "ADJUDICATED" | "APPROVED" | "REJECTED" | "MANUAL_REVIEW" | "PAID" | "VETTED" | "PARTIALLY_VETTED" | "READY_FOR_PAYMENT" | "PART_PAID";
export type ClaimType = "DIRECT" | "REIMBURSEMENT" | "CREDIT" | "EXCEPTION";
export type AdjudicationDecision = "APPROVE" | "REJECT" | "MANUAL_REVIEW";
export type RuleCategory = "eligibility" | "coverage" | "limits" | "fraud";
export type RuleResultStatus = "PASS" | "FAIL" | "FLAG";
export type FraudFlagType = "DUPLICATE" | "FREQUENCY" | "AMOUNT_THRESHOLD" | "EXPIRED_CONTRACT" | "SUSPENDED_PROVIDER" | "REPEAT_VISIT" | "RATE_CARD_OVERCHARGE";
export type FraudSeverity = "LOW" | "MEDIUM" | "HIGH" | "CRITICAL";
export type CaseStatus = "SCHEDULED" | "ADMITTED" | "IN_TREATMENT" | "DISCHARGED" | "CLOSED";

// ── Pre-Authorization ───────────────────────────────────────
export type PreAuthStatus = "SUBMITTED" | "UNDER_REVIEW" | "APPROVED" | "DENIED" | "INFO_REQUESTED" | "EXPIRED" | "CLAIMED";

// ── Provider Management ─────────────────────────────────────
export type ProviderStatus = "PENDING" | "CREDENTIALING" | "ACTIVE" | "SUSPENDED" | "TERMINATED";
export type ProviderType = "hospital" | "clinic" | "pharmacy" | "lab";
export type ProviderTier = "TIER_1" | "TIER_2" | "TIER_3";
export type ProviderNetworkStatus = "ACTIVE" | "INACTIVE";
export type AccreditationStatus = "NONE" | "PENDING" | "ACCREDITED" | "EXPIRED" | "REVOKED";
export type ContractStatus = "ACTIVE" | "EXPIRED" | "TERMINATED";

// ── Billing & Payments ──────────────────────────────────────
export type InvoiceStatus = "PENDING" | "PAID" | "OVERDUE" | "CANCELLED";
export type PaymentStatus = "INITIATED" | "PROCESSING" | "CONFIRMED" | "FAILED" | "RECONCILED" | "CANCELLED";
export type PaymentMethod = "MPESA" | "BANK_TRANSFER";
export type PaymentType = "PREMIUM" | "REMITTANCE";
export type RemittanceStatus = "PENDING" | "PROCESSING" | "SENT" | "CONFIRMED" | "FAILED";
export type InstallmentScheduleStatus = "ACTIVE" | "COMPLETED" | "CANCELLED";
export type InstallmentStatus = "PENDING" | "PAID" | "OVERDUE";
export type StatementStatus = "UPLOADED" | "RECONCILED";
export type MatchStatus = "UNMATCHED" | "MATCHED" | "DISPUTED";
export type Currency = "KES";

// ── Sales ───────────────────────────────────────────────────
export type LeadStatus = "NEW" | "CONTACTED" | "QUALIFIED" | "PROPOSAL_SENT" | "NEGOTIATION" | "WON" | "LOST" | "DORMANT";
export type LeadSource = "direct" | "referral" | "web" | "agent" | "broker";
export type LeadActivityType = "call" | "email" | "meeting" | "note" | "follow_up" | "status_change";
export type QuotationStatus = "DRAFT" | "ISSUED" | "PENDING_DECISION" | "ACCEPTED" | "DECLINED" | "EXPIRED" | "CONVERTED";
export type QuotationType = "standard" | "tailor_made";
export type ApprovalStatus = "NONE" | "PENDING" | "APPROVED" | "REJECTED";

// ── Notifications ───────────────────────────────────────────
export type NotificationChannel = "SMS" | "EMAIL" | "IN_APP" | "PUSH";
export type NotificationType = "QUOTATION" | "APPROVAL" | "CLAIM" | "POLICY" | "DOCUMENT";
export type NotificationStatus = "PENDING" | "SENT" | "DELIVERED" | "FAILED" | "READ";

// ── Audit Trail ─────────────────────────────────────────────
export type AuditAction = "CREATE" | "UPDATE" | "DELETE" | "STATE_CHANGE";
export type AuditEntityType =
  | "CLAIM" | "POLICY" | "MEMBER" | "PLAN" | "BENEFIT" | "EXCLUSION"
  | "PREMIUM_RULE" | "PROVIDER_NETWORK" | "PROVIDER" | "USER"
  | "LEAD" | "QUOTATION" | "QUOTATION_VERSION" | "QUOTATION_DOCUMENT"
  | "APPROVAL_LIMIT" | "ENDORSEMENT" | "RENEWAL" | "UNDERWRITING"
  | "POLICY_DOCUMENT" | "UNDERWRITING_FLAG" | "UNDERWRITING_RULE"
  | "CREDIT_NOTE" | "CASE_RECORD" | "CLAIM_DOCUMENT" | "PROVIDER_STATEMENT"
  | "TREATY" | "TREATY_PARTICIPANT" | "TREATY_LAYER" | "CESSION"
  | "REINSURANCE_RECOVERY" | "RECOVERY_WORKFLOW_EVENT" | "BORDEREAU"
  | "BORDEREAU_ITEM" | "REINSURER_STATEMENT" | "PROFIT_COMMISSION" | "TREATY_ALERT";

// ── Reinsurance ─────────────────────────────────────────────
export type TreatyType = "QUOTA_SHARE" | "XOL";
export type TreatyStatus = "DRAFT" | "ACTIVE" | "EXPIRED" | "TERMINATED";
export type CessionType = "PREMIUM" | "CLAIM";
export type CessionStatus = "PENDING" | "BOOKED" | "REVERSED";
export type RecoveryStatus = "NOTIFIED" | "ACKNOWLEDGED" | "INFO_REQUESTED" | "APPROVED" | "PAID" | "WRITTEN_OFF";
export type RecoveryWorkflowStatus = "NOTIFICATION" | "ACKNOWLEDGMENT" | "INFO_REQUEST" | "APPROVAL" | "PAYMENT";
export type BordereauType = "PREMIUM" | "CLAIM";
export type BordereauStatus = "DRAFT" | "FINALIZED" | "SENT";
export type ReinsurerStatementStatus = "DRAFT" | "ISSUED" | "ACKNOWLEDGED" | "SETTLED";
export type ProfitCommissionType = "SLIDING_SCALE" | "FLAT" | "CARRY_FORWARD";
export type TreatyAlertType = "LIMIT_BREACH" | "AGGREGATE_WARNING" | "CATASTROPHE_THRESHOLD" | "EXPIRY_WARNING";
export type TreatyAlertSeverity = "LOW" | "MEDIUM" | "HIGH" | "CRITICAL";
```

### Status Badge Color Mapping (Complete)

```typescript
// lib/utils/status-colors.ts — Color mapping for ALL status types

export const statusColors: Record<string, string> = {
  // Policy
  DRAFT: "bg-gray-100 text-gray-800",
  ACTIVE: "bg-green-100 text-green-800",
  LAPSED: "bg-yellow-100 text-yellow-800",
  TERMINATED: "bg-red-100 text-red-800",
  SUSPENDED: "bg-orange-100 text-orange-800",

  // Claims
  RECEIVED: "bg-blue-100 text-blue-800",
  VALIDATED: "bg-indigo-100 text-indigo-800",
  ADJUDICATED: "bg-purple-100 text-purple-800",
  APPROVED: "bg-green-100 text-green-800",
  REJECTED: "bg-red-100 text-red-800",
  MANUAL_REVIEW: "bg-amber-100 text-amber-800",
  PAID: "bg-emerald-100 text-emerald-800",
  PART_PAID: "bg-teal-100 text-teal-800",
  VETTED: "bg-cyan-100 text-cyan-800",
  PARTIALLY_VETTED: "bg-sky-100 text-sky-800",
  READY_FOR_PAYMENT: "bg-lime-100 text-lime-800",

  // Pre-Auth
  SUBMITTED: "bg-blue-100 text-blue-800",
  UNDER_REVIEW: "bg-purple-100 text-purple-800",
  DENIED: "bg-red-100 text-red-800",
  INFO_REQUESTED: "bg-amber-100 text-amber-800",
  EXPIRED: "bg-gray-100 text-gray-800",
  CLAIMED: "bg-green-100 text-green-800",

  // Provider
  PENDING: "bg-yellow-100 text-yellow-800",
  CREDENTIALING: "bg-indigo-100 text-indigo-800",

  // Underwriting
  DECLINED: "bg-red-100 text-red-800",
  REFER: "bg-amber-100 text-amber-800",

  // UW Flags
  OPEN: "bg-red-100 text-red-800",
  ACKNOWLEDGED: "bg-yellow-100 text-yellow-800",
  RESOLVED: "bg-green-100 text-green-800",
  OVERRIDDEN: "bg-purple-100 text-purple-800",

  // Case
  SCHEDULED: "bg-blue-100 text-blue-800",
  ADMITTED: "bg-indigo-100 text-indigo-800",
  IN_TREATMENT: "bg-purple-100 text-purple-800",
  DISCHARGED: "bg-teal-100 text-teal-800",
  CLOSED: "bg-gray-100 text-gray-800",

  // Billing
  OVERDUE: "bg-red-100 text-red-800",
  CANCELLED: "bg-gray-100 text-gray-800",

  // Payment
  INITIATED: "bg-blue-100 text-blue-800",
  PROCESSING: "bg-indigo-100 text-indigo-800",
  CONFIRMED: "bg-green-100 text-green-800",
  FAILED: "bg-red-100 text-red-800",
  RECONCILED: "bg-emerald-100 text-emerald-800",

  // Remittance
  SENT: "bg-teal-100 text-teal-800",

  // Statement
  UPLOADED: "bg-blue-100 text-blue-800",
  MATCHED: "bg-green-100 text-green-800",
  UNMATCHED: "bg-red-100 text-red-800",
  DISPUTED: "bg-amber-100 text-amber-800",

  // Credit Note
  APPLIED: "bg-green-100 text-green-800",

  // Lead
  NEW: "bg-blue-100 text-blue-800",
  CONTACTED: "bg-indigo-100 text-indigo-800",
  QUALIFIED: "bg-purple-100 text-purple-800",
  PROPOSAL_SENT: "bg-cyan-100 text-cyan-800",
  NEGOTIATION: "bg-amber-100 text-amber-800",
  WON: "bg-green-100 text-green-800",
  LOST: "bg-red-100 text-red-800",
  DORMANT: "bg-gray-100 text-gray-800",

  // Quotation
  ISSUED: "bg-blue-100 text-blue-800",
  PENDING_DECISION: "bg-amber-100 text-amber-800",
  ACCEPTED: "bg-green-100 text-green-800",
  DECLINED: "bg-red-100 text-red-800",
  CONVERTED: "bg-emerald-100 text-emerald-800",

  // Notification
  DELIVERED: "bg-green-100 text-green-800",
  READ: "bg-gray-100 text-gray-800",

  // Reinsurance
  BOOKED: "bg-green-100 text-green-800",
  REVERSED: "bg-red-100 text-red-800",
  NOTIFIED: "bg-blue-100 text-blue-800",
  WRITTEN_OFF: "bg-gray-100 text-gray-800",
  FINALIZED: "bg-indigo-100 text-indigo-800",
  SETTLED: "bg-emerald-100 text-emerald-800",

  // Accreditation
  ACCREDITED: "bg-green-100 text-green-800",
  REVOKED: "bg-red-100 text-red-800",
  NONE: "bg-gray-100 text-gray-800",

  // Severity (used as secondary badge)
  LOW: "bg-blue-100 text-blue-800",
  MEDIUM: "bg-yellow-100 text-yellow-800",
  HIGH: "bg-orange-100 text-orange-800",
  CRITICAL: "bg-red-100 text-red-800",

  // Provider Tier
  TIER_1: "bg-green-100 text-green-800",
  TIER_2: "bg-blue-100 text-blue-800",
  TIER_3: "bg-gray-100 text-gray-800",
};

export function getStatusColor(status: string): string {
  return statusColors[status] || "bg-gray-100 text-gray-800";
}
```

### Enum Display Labels

```typescript
// lib/utils/enum-labels.ts — Human-readable labels for enum values

export const CLAIM_STATUS_LABELS: Record<ClaimStatus, string> = {
  RECEIVED: "Received",
  VALIDATED: "Validated",
  ADJUDICATED: "Adjudicated",
  APPROVED: "Approved",
  REJECTED: "Rejected",
  MANUAL_REVIEW: "Manual Review",
  PAID: "Paid",
  PART_PAID: "Partially Paid",
  VETTED: "Vetted",
  PARTIALLY_VETTED: "Partially Vetted",
  READY_FOR_PAYMENT: "Ready for Payment",
};

export const POLICY_STATUS_LABELS: Record<PolicyStatus, string> = {
  DRAFT: "Draft",
  ACTIVE: "Active",
  LAPSED: "Lapsed",
  TERMINATED: "Terminated",
  SUSPENDED: "Suspended",
};

export const PROVIDER_STATUS_LABELS: Record<ProviderStatus, string> = {
  PENDING: "Pending",
  CREDENTIALING: "Credentialing",
  ACTIVE: "Active",
  SUSPENDED: "Suspended",
  TERMINATED: "Terminated",
};

export const CASE_STATUS_LABELS: Record<CaseStatus, string> = {
  SCHEDULED: "Scheduled",
  ADMITTED: "Admitted",
  IN_TREATMENT: "In Treatment",
  DISCHARGED: "Discharged",
  CLOSED: "Closed",
};

export const PREAUTH_STATUS_LABELS: Record<PreAuthStatus, string> = {
  SUBMITTED: "Submitted",
  UNDER_REVIEW: "Under Review",
  APPROVED: "Approved",
  DENIED: "Denied",
  INFO_REQUESTED: "Info Requested",
  EXPIRED: "Expired",
  CLAIMED: "Claimed",
};

export const LEAD_STATUS_LABELS: Record<LeadStatus, string> = {
  NEW: "New",
  CONTACTED: "Contacted",
  QUALIFIED: "Qualified",
  PROPOSAL_SENT: "Proposal Sent",
  NEGOTIATION: "Negotiation",
  WON: "Won",
  LOST: "Lost",
  DORMANT: "Dormant",
};

export const RECOVERY_STATUS_LABELS: Record<RecoveryStatus, string> = {
  NOTIFIED: "Notified",
  ACKNOWLEDGED: "Acknowledged",
  INFO_REQUESTED: "Info Requested",
  APPROVED: "Approved",
  PAID: "Paid",
  WRITTEN_OFF: "Written Off",
};

export const UNDERWRITING_STATUS_LABELS: Record<UnderwritingStatus, string> = {
  PENDING: "Pending",
  APPROVED: "Approved",
  DECLINED: "Declined",
  REFER: "Referred",
};

export const PAYMENT_STATUS_LABELS: Record<PaymentStatus, string> = {
  INITIATED: "Initiated",
  PROCESSING: "Processing",
  CONFIRMED: "Confirmed",
  FAILED: "Failed",
  RECONCILED: "Reconciled",
  CANCELLED: "Cancelled",
};

export const ACCREDITATION_LABELS: Record<AccreditationStatus, string> = {
  NONE: "None",
  PENDING: "Pending",
  ACCREDITED: "Accredited",
  EXPIRED: "Expired",
  REVOKED: "Revoked",
};

export const UW_FLAG_STATUS_LABELS: Record<UnderwritingFlagStatus, string> = {
  OPEN: "Open",
  ACKNOWLEDGED: "Acknowledged",
  RESOLVED: "Resolved",
  OVERRIDDEN: "Overridden",
};

export const TREATY_STATUS_LABELS: Record<TreatyStatus, string> = {
  DRAFT: "Draft",
  ACTIVE: "Active",
  EXPIRED: "Expired",
  TERMINATED: "Terminated",
};

export const BORDEREAU_STATUS_LABELS: Record<BordereauStatus, string> = {
  DRAFT: "Draft",
  FINALIZED: "Finalized",
  SENT: "Sent",
};

export const REINSURER_STATEMENT_LABELS: Record<ReinsurerStatementStatus, string> = {
  DRAFT: "Draft",
  ISSUED: "Issued",
  ACKNOWLEDGED: "Acknowledged",
  SETTLED: "Settled",
};
```

### Auth Types

```typescript
// lib/types/auth.ts

interface LoginRequest {
  email: string;
  password: string;
}

interface RegisterRequest {
  email: string;
  password: string;
  name: string;
  phone: string;
  national_id?: string;
  role_name?: string;
}

interface RefreshTokenRequest {
  refresh_token: string;
}

interface LoginResponse {
  access_token: string;
  access_token_expires_at: ISODateTime;
  refresh_token: string;
  user: UserResponse;
}

interface RegisterResponse {
  user_id: string;
  email: string;
  message: string;
}

interface UserResponse {
  id: UUID;
  email: string;
  name: string;
  phone: string;
  national_id: string;
  role_id: UUID;
  role_name: string;
  status: "ACTIVE" | "INACTIVE" | "SUSPENDED";
  created_at: ISODateTime;
  updated_at: ISODateTime;
}

interface CreateUserRequest {
  email: string;        // required, valid email
  name: string;         // required
  phone: string;        // required
  national_id?: string;
  role_name: string;    // required
  password: string;     // required, min 8 chars
}

interface UpdateUserRequest {
  name?: string;
  phone?: string;
  national_id?: string;
}

interface AssignRoleRequest {
  role_id: UUID;  // required
}

interface UpdateUserStatusRequest {
  status: "ACTIVE" | "INACTIVE" | "SUSPENDED";  // required
}

type UserRole = "Admin" | "Underwriter" | "ClaimsOfficer" | "Finance" | "Provider" | "Member" | "SalesAgent" | "Manager";
```

### Product Types

```typescript
// lib/types/product.ts

interface PlanResponse {
  id: UUID;
  name: string;
  type: "individual" | "group";
  segment: "retail" | "corporate" | "sme";
  base_premium: Money;
  currency: string;
  status: "ACTIVE" | "INACTIVE";
  description: string;
  created_at: ISODateTime;
  updated_at: ISODateTime;
}

interface CreatePlanRequest {
  name: string;
  type: "individual" | "group";
  segment?: "retail" | "corporate" | "sme";
  base_premium: Money;
  currency?: string;
  description?: string;
}

interface BenefitResponse {
  id: UUID;
  plan_id: UUID;
  parent_benefit_id: UUID | null;
  name: string;
  category: "outpatient" | "inpatient" | "dental" | "optical" | "maternity";
  annual_limit: Money;
  co_pay_type: "percentage" | "fixed";
  co_pay_value: Money;
  waiting_period_days: number;
  sub_limit_type: "none" | "per_visit" | "per_item";
  sub_limit_value: Money;
  min_age: number;
  max_age: number;
  waiting_period_type: "general" | "maternity" | "pre_existing" | "chronic" | "surgical";
  deductible_amount: Money;
  created_at: ISODateTime;
}

interface CreateBenefitRequest {
  parent_benefit_id?: UUID;
  name: string;
  category: string;
  annual_limit: Money;
  co_pay_type: "percentage" | "fixed";
  co_pay_value: Money;
  waiting_period_days?: number;
  sub_limit_type?: string;
  sub_limit_value?: Money;
  min_age?: number;
  max_age?: number;
  waiting_period_type?: string;
  deductible_amount?: Money;
}

interface ExclusionResponse {
  id: UUID;
  plan_id: UUID;
  description: string;
  type: "pre_existing" | "cosmetic" | "experimental";
  icd_codes: string[];
  created_at: ISODateTime;
}

interface CreateExclusionRequest {
  description: string;  // required
  type: "pre_existing" | "cosmetic" | "experimental";  // required
  icd_codes?: string[];
}

interface UpdateExclusionRequest {
  description?: string;
  type?: string;
  icd_codes?: string[];
}

interface UpdatePlanRequest {
  name?: string;
  type?: string;
  segment?: string;
  base_premium?: Money;
  description?: string;
  status?: string;
}

interface PremiumRuleResponse {
  id: UUID;
  plan_id: UUID;
  calculation_type: "per_member" | "per_family" | "tiered" | "flat";
  relationship: string;
  rate_amount: Money;
  discount_type: "percentage" | "fixed";
  discount_value: Money;
  min_members: number;
  min_age: number;
  max_age: number;
  created_at: ISODateTime;
  updated_at: ISODateTime;
}

interface UnderwritingRuleResponse {
  id: UUID;
  plan_id: UUID;
  rule_type: "MAX_AGE" | "MIN_AGE" | "DOUBLE_INSURANCE" | "PRE_EXISTING_CONDITION" | "BMI_THRESHOLD" | "WAITING_PERIOD";
  relationship: string;
  parameter_key: string;
  parameter_value: string;
  severity: "LOW" | "MEDIUM" | "HIGH";
  risk_score_weight: number;
  is_blocking: boolean;
  is_active: boolean;
  description: string;
  created_at: ISODateTime;
  updated_at: ISODateTime;
}

interface CreatePremiumRuleRequest {
  calculation_type: string;  // required
  relationship?: string;
  rate_amount: Money;        // required, min 1
  discount_type?: string;
  discount_value?: Money;
  min_members?: number;
  min_age?: number;
  max_age?: number;
}

interface UpdatePremiumRuleRequest {
  calculation_type?: string;
  relationship?: string;
  rate_amount?: Money;
  discount_type?: string;
  discount_value?: Money;
  min_members?: number;
  min_age?: number;
  max_age?: number;
}

interface CreateUnderwritingRuleRequest {
  plan_id: UUID;              // required
  rule_type: string;          // required
  relationship?: string;
  parameter_key: string;      // required
  parameter_value: string;    // required
  severity?: string;
  risk_score_weight?: number;
  is_blocking?: boolean;
  is_active?: boolean;
  description?: string;
}

interface UpdateUnderwritingRuleRequest {
  rule_type?: string;
  relationship?: string;
  parameter_key?: string;
  parameter_value?: string;
  severity?: string;
  risk_score_weight?: number;
  is_blocking?: boolean;
  is_active?: boolean;
  description?: string;
}

interface CreateProviderNetworkRequest {
  provider_id: UUID;  // required
  benefit_category?: string;
}

interface UpdateProviderNetworkStatusRequest {
  status: string;  // required
}

interface ProviderNetworkResponse {
  id: UUID;
  plan_id: UUID;
  provider_id: UUID;
  benefit_category: string;
  status: "ACTIVE" | "INACTIVE";
  created_at: ISODateTime;
  updated_at: ISODateTime;
}
```

### Policy Types

```typescript
// lib/types/policy.ts

type PolicyStatus = "DRAFT" | "ACTIVE" | "LAPSED" | "TERMINATED" | "SUSPENDED";
type MemberStatus = "ACTIVE" | "SUSPENDED" | "REMOVED";

interface PolicyResponse {
  id: UUID;
  plan_id: UUID;
  policyholder_name: string;
  policyholder_email: string;
  policyholder_phone: string;
  policy_number: string;
  status: PolicyStatus;
  start_date: ISODateTime;
  end_date: ISODateTime;
  premium_amount: Money;
  currency: string;
  renewed_from_id?: UUID;
  created_at: ISODateTime;
  updated_at: ISODateTime;
}

interface CreatePolicyRequest {
  plan_id: UUID;
  policyholder_name: string;
  policyholder_email: string;
  policyholder_phone: string;
  start_date?: ISODateTime;
  end_date?: ISODateTime;
}

interface MemberResponse {
  id: UUID;
  policy_id: UUID;
  national_id: string;
  name: string;
  date_of_birth: ISODateTime;
  gender: "male" | "female" | "other";
  relationship: "principal" | "spouse" | "child" | "parent";
  member_number: string;
  phone: string;
  email: string;
  kra_pin: string;
  county: string;
  address: string;
  status: MemberStatus;
  verified: boolean;
  verified_at?: ISODateTime;
  created_at: ISODateTime;
}

interface EnrollMemberRequest {
  national_id?: string;
  name: string;
  date_of_birth: string; // YYYY-MM-DD
  gender: "male" | "female" | "other";
  relationship: "principal" | "spouse" | "child" | "parent";
  phone?: string;
  email?: string;
  kra_pin?: string;
  county?: string;
  address?: string;
}

interface ActivatePolicyRequest {
  payment_reference: string;  // required
}

interface UpdatePolicyRequest {
  policyholder_name?: string;
  policyholder_email?: string;
  policyholder_phone?: string;
  start_date?: ISODateTime;
  end_date?: ISODateTime;
}

interface ChangePlanRequest {
  new_plan_id: UUID;  // required
  reason?: string;
}

interface UpdateMemberRequest {
  name?: string;
  phone?: string;
  email?: string;
  kra_pin?: string;
  county?: string;
  address?: string;
}

interface RemoveMemberRequest {
  reason?: string;
}

interface BulkIDsRequest {
  ids: string[];  // required — used for bulk activate/lapse
}

interface BulkEnrollRequest {
  members: EnrollMemberRequest[];  // required
}

interface BulkRemoveRequest {
  member_ids: string[];  // required
  reason?: string;
}

interface BulkMemberResultResponse {
  succeeded: number;
  failed: number;
  members?: MemberResponse[];
  errors?: string[];
}

interface EndorsementResponse {
  id: UUID;
  policy_id: UUID;
  endorsement_type: "ADD_MEMBER" | "REMOVE_MEMBER" | "UPDATE_MEMBER" | "PLAN_CHANGE";
  status: "PENDING" | "APPROVED" | "REJECTED" | "APPLIED";
  effective_date: ISODateTime;
  changes: Record<string, unknown>;
  reason: string;
  premium_adjustment: Money;
  requested_by: UUID;
  approved_by?: UUID;
  approved_at?: ISODateTime;
  applied_at?: ISODateTime;
  created_at: ISODateTime;
  updated_at: ISODateTime;
}

interface RenewalResponse {
  id: UUID;
  policy_id: UUID;
  renewed_policy_id?: UUID;
  status: "PENDING" | "APPROVED" | "REJECTED" | "COMPLETED" | "EXPIRED";
  renewal_date: ISODateTime;
  new_premium: Money;
  premium_change_reason: string;
  new_plan_id?: UUID;
  approved_by?: UUID;
  approved_at?: ISODateTime;
  completed_at?: ISODateTime;
  expires_at?: ISODateTime;
  created_by: UUID;
  created_at: ISODateTime;
  updated_at: ISODateTime;
}

interface UnderwritingResponse {
  id: UUID;
  policy_id: UUID;
  member_id?: UUID;
  status: "PENDING" | "APPROVED" | "DECLINED" | "REFER";
  questionnaire: Record<string, unknown>;
  medical_declarations: Record<string, unknown>;
  risk_score: number;
  risk_flags: string[];
  decision_reason: string;
  assessed_by?: UUID;
  assessed_at?: ISODateTime;
  created_by: UUID;
  created_at: ISODateTime;
  updated_at: ISODateTime;
}

interface UnderwritingFlagResponse {
  id: UUID;
  assessment_id?: UUID;
  policy_id: UUID;
  member_id?: UUID;
  flag_type: "MAX_AGE" | "MIN_AGE" | "DOUBLE_INSURANCE" | "PRE_EXISTING_CONDITION" | "BMI_THRESHOLD" | "WAITING_PERIOD" | "RENEWAL_SKIP";
  severity: "LOW" | "MEDIUM" | "HIGH";
  details: string;
  status: "OPEN" | "ACKNOWLEDGED" | "RESOLVED" | "OVERRIDDEN";
  resolved_by?: UUID;
  resolved_at?: ISODateTime;
  resolution?: string;
  created_at: ISODateTime;
  updated_at: ISODateTime;
}

interface PolicyDocumentResponse {
  id: UUID;
  policy_id: UUID;
  member_id?: UUID;
  document_type: "WELCOME_LETTER" | "MEMBER_CARD" | "POLICY_SCHEDULE" | "RENEWAL_NOTICE" | "ENDORSEMENT" | "LOU" | "DECLINE_LETTER";
  file_name: string;
  file_size: number;
  s3_key: string;
  generated_by: UUID;
  created_at: ISODateTime;
}

interface CreateEndorsementRequest {
  policy_id: UUID;                     // required
  endorsement_type: "ADD_MEMBER" | "REMOVE_MEMBER" | "UPDATE_MEMBER" | "PLAN_CHANGE";  // required
  effective_date: string;              // required, YYYY-MM-DD
  changes: Record<string, unknown>;    // required — JSON object with change details
  reason?: string;
}

interface RejectEndorsementRequest {
  reason: string;  // required
}

interface InitiateRenewalRequest {
  policy_id: UUID;       // required
  new_plan_id?: UUID;
  renewal_date: string;  // required, YYYY-MM-DD
  expires_at?: string;
}

interface RejectRenewalRequest {
  reason: string;  // required
}

interface BulkRenewalRequest {
  policy_ids: string[];  // required
}

interface SubmitAssessmentRequest {
  policy_id: UUID;                              // required
  member_id?: UUID;
  questionnaire: Record<string, unknown>;       // required
  medical_declarations?: Record<string, unknown>;
}

interface ReviewAssessmentRequest {
  status: string;  // required — "APPROVED" | "DECLINED" | "REFER"
  risk_score?: number;
  decision_reason?: string;
}

interface ResolveFlagRequest {
  resolution: string;  // required
}

interface OverrideFlagRequest {
  reason: string;  // required
}

interface ApplyCreditNoteRequest {
  invoice_id: UUID;  // required
}

interface CreditNoteResponse {
  id: UUID;
  policy_id: UUID;
  member_id?: UUID;
  credit_note_number: string;
  amount: Money;
  currency: string;
  reason: string;
  status: "DRAFT" | "APPROVED" | "APPLIED" | "CANCELLED";
  applied_to_invoice_id?: UUID;
  approved_by?: UUID;
  approved_at?: ISODateTime;
  applied_at?: ISODateTime;
  created_by: UUID;
  created_at: ISODateTime;
  updated_at: ISODateTime;
}
```

### Claims Types

```typescript
// lib/types/claims.ts

type ClaimStatus = "RECEIVED" | "VALIDATED" | "ADJUDICATED" | "APPROVED" | "REJECTED" | "MANUAL_REVIEW" | "PAID" | "VETTED" | "PARTIALLY_VETTED" | "READY_FOR_PAYMENT" | "PART_PAID";

interface ClaimResponse {
  id: UUID;
  claim_number: string;
  policy_id: UUID;
  member_id: UUID;
  provider_id: UUID;
  status: ClaimStatus;
  total_amount: Money;
  approved_amount: Money;
  co_pay_amount: Money;
  member_responsibility: Money;
  diagnosis_codes: string[];
  service_date: ISODateTime;
  notes: string;
  claim_type: "DIRECT" | "REIMBURSEMENT" | "CREDIT" | "EXCEPTION";
  vetted_amount?: Money;
  vetted_by?: UUID;
  vetted_at?: ISODateTime;
  sla_breach_at?: ISODateTime;
  rejection_reason?: string;
  line_items?: LineItemResponse[];
  decision?: AdjudicationResponse;
  fraud_flags?: FraudFlagResponse[];
  created_at: ISODateTime;
  updated_at: ISODateTime;
}

interface SubmitClaimRequest {
  policy_id: UUID;
  member_id: UUID;
  provider_id: UUID;
  pre_auth_id?: UUID;
  diagnosis_codes: string[];
  service_date: ISODateTime;
  admission_date?: ISODateTime;
  discharge_date?: ISODateTime;
  notes?: string;
  claim_type?: string;
  line_items: LineItemRequest[];
}

interface LineItemRequest {
  procedure_code: string;
  procedure_name: string;
  diagnosis_code?: string;
  quantity: number;
  unit_price: Money;
}

interface LineItemResponse {
  id: UUID;
  procedure_code: string;
  procedure_name: string;
  diagnosis_code: string;
  quantity: number;
  unit_price: Money;
  total_price: Money;
  approved_amount: Money;
}

interface AdjudicationResponse {
  decision: "APPROVE" | "REJECT" | "MANUAL_REVIEW";
  payable_amount: Money;
  member_responsibility: Money;
  deductible_applied: Money;
  co_pay_applied: Money;
  sub_limit_applied?: Money;
  benefit_category?: string;
  reasons: Record<string, unknown>;
  rule_results: Record<string, unknown>;
  adjudicated_at: ISODateTime;
}

interface FraudFlagResponse {
  id: UUID;
  flag_type: "DUPLICATE" | "FREQUENCY" | "AMOUNT_THRESHOLD" | "EXPIRED_CONTRACT" | "SUSPENDED_PROVIDER" | "REPEAT_VISIT" | "RATE_CARD_OVERCHARGE";
  severity: "LOW" | "MEDIUM" | "HIGH" | "CRITICAL";
  details: string;
  resolved: boolean;
}

interface ReviewClaimRequest {
  decision: "APPROVED" | "REJECTED";  // required
  reason?: string;
}

interface VetClaimRequest {
  vetted_amount: Money;  // required, min 0
  notes?: string;
}

interface BulkSubmitClaimsRequest {
  claims: SubmitClaimRequest[];  // required, min 1
}

interface BulkClaimResultResponse {
  succeeded: number;
  failed: number;
  claims?: ClaimResponse[];
  errors?: string[];
}

interface CreateCaseRequest {
  preauth_id: UUID;                  // required
  expected_discharge?: ISODateTime;
  diagnosis?: string;
  treating_doctor?: string;
  room_type?: string;
  estimated_cost?: Money;
  notes?: string;
}

interface UpdateCaseRequest {
  diagnosis?: string;
  treating_doctor?: string;
  room_type?: string;
  estimated_cost?: Money;
  notes?: string;
}

interface AdmitCaseRequest {
  admission_date: ISODateTime;  // required
}

interface DischargeCaseRequest {
  actual_discharge: ISODateTime;  // required
  actual_cost: Money;             // required, min 0
}

interface SubmitPreAuthRequest {
  policy_id: UUID;          // required
  member_id: UUID;          // required
  provider_id: UUID;        // required
  procedure_codes: string[];  // required
  diagnosis_codes: string[];  // required
  estimated_cost: Money;      // required, min 1
  notes?: string;
}

interface ReviewPreAuthRequest {
  decision: "APPROVED" | "DENIED" | "INFO_REQUESTED";  // required
  approved_amount?: Money;
  denial_reason?: string;
  validity_days?: number;
}

interface PreAuthResponse {
  id: UUID;
  policy_id: UUID;
  member_id: UUID;
  provider_id: UUID;
  auth_code: string;
  procedure_codes: string[];
  diagnosis_codes: string[];
  estimated_cost: Money;
  approved_amount: Money;
  status: "SUBMITTED" | "UNDER_REVIEW" | "APPROVED" | "DENIED" | "INFO_REQUESTED" | "EXPIRED" | "CLAIMED";
  validity_start?: ISODateTime;
  validity_end?: ISODateTime;
  notes: string;
  denial_reason?: string;
  created_at: ISODateTime;
  updated_at: ISODateTime;
}

interface CaseRecordResponse {
  id: UUID;
  case_number: string;
  pre_auth_id: UUID;
  policy_id: UUID;
  member_id: UUID;
  provider_id: UUID;
  status: "SCHEDULED" | "ADMITTED" | "IN_TREATMENT" | "DISCHARGED" | "CLOSED";
  admission_date?: ISODateTime;
  expected_discharge?: ISODateTime;
  actual_discharge?: ISODateTime;
  diagnosis: string;
  treating_doctor: string;
  room_type: string;
  total_estimated_cost: Money;
  total_actual_cost: Money;
  notes: string;
  closed_at?: ISODateTime;
  created_by: UUID;
  created_at: ISODateTime;
  updated_at: ISODateTime;
}

interface ClaimDocumentResponse {
  id: UUID;
  claim_id: UUID;
  file_name: string;
  file_type: string;
  file_size: number;
  s3_key: string;
  uploaded_by: UUID;
  created_at: ISODateTime;
}
```

### Sales Types

```typescript
// lib/types/sales.ts

type LeadStatus = "NEW" | "CONTACTED" | "QUALIFIED" | "PROPOSAL_SENT" | "NEGOTIATION" | "WON" | "LOST" | "DORMANT";

interface CreateLeadRequest {
  contact_name: string;    // required
  contact_email?: string;
  contact_phone?: string;
  company_name?: string;
  source: "direct" | "referral" | "web" | "agent" | "broker";  // required
  segment: "retail" | "corporate" | "sme";                      // required
  plan_type: "individual" | "group";                             // required
  estimated_members?: number;
  expected_premium?: Money;
  closure_probability?: number;
  next_follow_up_date?: ISODateTime;
  notes?: string;
}

interface UpdateLeadRequest {
  contact_name?: string;
  contact_email?: string;
  contact_phone?: string;
  company_name?: string;
  source?: string;
  segment?: string;
  plan_type?: string;
  estimated_members?: number;
  expected_premium?: Money;
  closure_probability?: number;
  assigned_to?: string;
  next_follow_up_date?: ISODateTime;
  notes?: string;
}

interface UpdateLeadStatusRequest {
  status: LeadStatus;  // required
}

interface CreateLeadActivityRequest {
  activity_type: "call" | "email" | "meeting" | "note" | "follow_up";  // required
  description?: string;
  scheduled_at?: ISODateTime;
  completed_at?: ISODateTime;
}

interface LeadResponse {
  id: UUID;
  lead_number: string;
  contact_name: string;
  contact_email: string;
  contact_phone: string;
  company_name: string;
  source: "direct" | "referral" | "web" | "agent" | "broker";
  segment: "retail" | "corporate" | "sme";
  plan_type: "individual" | "group";
  estimated_members: number;
  expected_premium: Money;
  closure_probability: number;
  currency: string;
  status: LeadStatus;
  assigned_to: UUID;
  next_follow_up_date?: ISODateTime;
  notes: string;
  created_by: UUID;
  created_at: ISODateTime;
  updated_at: ISODateTime;
}

interface LeadActivityResponse {
  id: UUID;
  lead_id: UUID;
  activity_type: "call" | "email" | "meeting" | "note" | "follow_up";
  description: string;
  scheduled_at?: ISODateTime;
  completed_at?: ISODateTime;
  created_by: UUID;
  created_at: ISODateTime;
}

interface QuotationResponse {
  id: UUID;
  quotation_number: string;
  lead_id: UUID;
  plan_id: UUID;
  quotation_type: "standard" | "tailor_made";
  status: "DRAFT" | "ISSUED" | "PENDING_DECISION" | "ACCEPTED" | "DECLINED" | "EXPIRED" | "CONVERTED";
  current_version: number;
  policy_id?: UUID;
  valid_from?: ISODateTime;
  valid_until?: ISODateTime;
  client_name: string;
  client_email: string;
  client_phone: string;
  currency: string;
  created_by: UUID;
  created_at: ISODateTime;
  updated_at: ISODateTime;
}

interface QuotationVersionResponse {
  id: UUID;
  quotation_id: UUID;
  version_number: number;
  base_premium: Money;
  discount_type: string;
  discount_value: Money;
  discount_reason: string;
  loading_type: string;
  loading_value: Money;
  loading_reason: string;
  final_premium: Money;
  member_count: number;
  proposed_members: Record<string, unknown>[];
  billing_frequency: "monthly" | "quarterly" | "semi_annual" | "annual";
  requires_approval: boolean;
  approval_status: "NONE" | "PENDING" | "APPROVED" | "REJECTED";
  approved_by?: UUID;
  approved_at?: ISODateTime;
  rejection_reason?: string;
  pricing_breakdown: Record<string, unknown>;
  created_by: UUID;
  created_at: ISODateTime;
  updated_at: ISODateTime;
}

interface QuotationDocumentResponse {
  id: UUID;
  quotation_id: UUID;
  version_number: number;
  file_name: string;
  file_type: string;
  file_size: number;
  uploaded_by: UUID;
  can_edit_roles: string[];
  can_delete_roles: string[];
  created_at: ISODateTime;
}

interface ApprovalLimitResponse {
  id: UUID;
  role_name: string;
  max_discount_percentage: Money;
  max_discount_amount: Money;
  max_loading_percentage: Money;
  max_loading_amount: Money;
  escalation_role: string;
  is_active: boolean;
  created_at: ISODateTime;
  updated_at: ISODateTime;
}

interface VersionComparisonResponse {
  version_a: QuotationVersionResponse;
  version_b: QuotationVersionResponse;
  pricing_diff: PricingDiff;
}

interface PricingDiff {
  base_premium_diff: Money;
  discount_diff: Money;
  loading_diff: Money;
  final_premium_diff: Money;
  member_count_diff: number;
}

interface CreateQuotationRequest {
  lead_id: UUID;            // required
  plan_id: UUID;            // required
  quotation_type: "standard" | "tailor_made";  // required
  client_name: string;      // required
  client_email?: string;
  client_phone?: string;
  member_count: number;     // required, min 1
  proposed_members?: Record<string, unknown>[];
  billing_frequency: "monthly" | "quarterly" | "semi_annual" | "annual";  // required
  discount_type?: string;
  discount_value?: Money;
  discount_reason?: string;
  loading_type?: string;
  loading_value?: Money;
  loading_reason?: string;
}

interface CreateQuotationVersionRequest {
  member_count: number;     // required, min 1
  proposed_members?: Record<string, unknown>[];
  billing_frequency: "monthly" | "quarterly" | "semi_annual" | "annual";  // required
  discount_type?: string;
  discount_value?: Money;
  discount_reason?: string;
  loading_type?: string;
  loading_value?: Money;
  loading_reason?: string;
}

interface ApproveVersionRequest {
  notes?: string;
}

interface RejectVersionRequest {
  reason: string;  // required
}

interface SendQuotationRequest {
  channel: "SMS" | "EMAIL";  // required
  message?: string;
}

interface ConvertToPolicyRequest {
  start_date: string;  // required, YYYY-MM-DD
  notes?: string;
}

interface UploadDocumentMeta {
  file_name: string;    // required
  file_type: string;    // required
  file_size: number;    // required
  version_number?: number;
  can_edit_roles?: string[];
  can_delete_roles?: string[];
}

interface UpdateDocumentMeta {
  file_name?: string;
  can_edit_roles?: string[];
  can_delete_roles?: string[];
}

interface CreateApprovalLimitRequest {
  role_name: string;  // required
  max_discount_percentage?: Money;
  max_discount_amount?: Money;
  max_loading_percentage?: Money;
  max_loading_amount?: Money;
  escalation_role?: string;
}

interface UpdateApprovalLimitRequest {
  max_discount_percentage?: Money;
  max_discount_amount?: Money;
  max_loading_percentage?: Money;
  max_loading_amount?: Money;
  escalation_role?: string;
}

interface QuotationDetailResponse extends QuotationResponse {
  versions: QuotationVersionResponse[];
  documents: QuotationDocumentResponse[];
}

interface ConversionResultResponse {
  quotation_id: UUID;
  policy_id: UUID;
  quotation_number: string;
  policy_number: string;
  message: string;
}
```

### Billing Types

```typescript
// lib/types/billing.ts

interface InvoiceResponse {
  id: UUID;
  policy_id: UUID;
  invoice_number: string;
  amount: Money;
  currency: string;
  due_date: ISODateTime;
  status: "PENDING" | "PAID" | "OVERDUE" | "CANCELLED";
  billing_period_start: ISODateTime;
  billing_period_end: ISODateTime;
  created_at: ISODateTime;
}

interface PaymentResponse {
  id: UUID;
  type: "PREMIUM" | "REMITTANCE";
  amount: Money;
  currency: string;
  method: "MPESA" | "BANK_TRANSFER";
  reference_number: string;
  status: "INITIATED" | "PROCESSING" | "CONFIRMED" | "FAILED" | "RECONCILED" | "CANCELLED";
  retry_count: number;
  paid_at?: ISODateTime;
  created_at: ISODateTime;
}

interface RemittanceResponse {
  id: UUID;
  provider_id: UUID;
  total_amount: Money;
  currency: string;
  status: "PENDING" | "PROCESSING" | "SENT" | "CONFIRMED" | "FAILED";
  remittance_advice_sent: boolean;
  period_start: ISODateTime;
  period_end: ISODateTime;
  created_at: ISODateTime;
}

interface InstallmentScheduleResponse {
  id: UUID;
  policy_id: UUID;
  frequency: "monthly" | "quarterly" | "semi_annual" | "annual";
  total_installments: number;
  amount_per_installment: Money;
  start_date: ISODateTime;
  status: "ACTIVE" | "COMPLETED" | "CANCELLED";
  created_at: ISODateTime;
  installments?: InstallmentResponse[];
}

interface InstallmentResponse {
  id: UUID;
  schedule_id: UUID;
  installment_number: number;
  due_date: ISODateTime;
  amount: Money;
  status: "PENDING" | "PAID" | "OVERDUE";
  paid_at?: ISODateTime;
  invoice_id?: UUID;
  created_at: ISODateTime;
}

interface ProviderStatementResponse {
  id: UUID;
  provider_id: UUID;
  statement_number: string;
  period_start: ISODateTime;
  period_end: ISODateTime;
  total_claimed: Money;
  total_matched: Money;
  total_discrepancy: Money;
  matched_count: number;
  unmatched_count: number;
  status: "UPLOADED" | "RECONCILED";
  file_name?: string;
  reconciled_at?: ISODateTime;
  created_at: ISODateTime;
}

interface InitiatePaymentRequest {
  invoice_id?: UUID;
  claim_id?: UUID;
  amount: Money;   // required, min 1
  method: "MPESA" | "BANK_TRANSFER";  // required
  phone?: string;
}

interface CreateInstallmentScheduleRequest {
  policy_id: UUID;       // required
  frequency: string;     // required — "monthly" | "quarterly" | "semi_annual" | "annual"
  start_date?: ISODateTime;
}

interface MarkInstallmentPaidRequest {
  invoice_id?: UUID;
}

interface UploadStatementRequest {
  provider_id: UUID;          // required
  period_start: ISODateTime;  // required
  period_end: ISODateTime;    // required
  file_name?: string;
  s3_key?: string;
  line_items: StatementLineItemInput[];  // required, min 1
}

interface StatementLineItemInput {
  claim_number?: string;
  service_date?: ISODateTime;
  member_name?: string;
  procedure_code?: string;
  claimed_amount: Money;  // required, min 0
}

interface StatementLineItemResponse {
  id: UUID;
  statement_id: UUID;
  claim_number?: string;
  service_date?: ISODateTime;
  member_name?: string;
  procedure_code?: string;
  claimed_amount: Money;
  matched_claim_id?: UUID;
  match_status: string;
  discrepancy_amount: Money;
  notes?: string;
  created_at: ISODateTime;
}

interface PaymentExportResponse {
  remittance_id: UUID;
  provider_id: UUID;
  provider_name: string;
  total_amount: Money;
  currency: string;
  period_start: ISODateTime;
  period_end: ISODateTime;
  claims: PaymentExportClaim[];
}

interface PaymentExportClaim {
  claim_number: string;
  amount: Money;
  service_date: ISODateTime;
}
```

### Provider Types

```typescript
// lib/types/provider.ts

interface ProviderResponse {
  id: UUID;
  name: string;
  type: "hospital" | "clinic" | "pharmacy" | "lab";
  license_number: string;
  status: "PENDING" | "CREDENTIALING" | "ACTIVE" | "SUSPENDED" | "TERMINATED";
  tier: "TIER_1" | "TIER_2" | "TIER_3";
  county: string;
  phone: string;
  email: string;
  contact_person: string;
  accreditation_status: "NONE" | "PENDING" | "ACCREDITED" | "EXPIRED" | "REVOKED";
  accreditation_expiry?: ISODateTime;
  accreditation_body?: string;
  created_at: ISODateTime;
  updated_at: ISODateTime;
}

interface ContractResponse {
  id: UUID;
  provider_id: UUID;
  start_date: ISODateTime;
  end_date: ISODateTime;
  terms: string;
  status: "ACTIVE" | "EXPIRED" | "TERMINATED";
  created_at: ISODateTime;
}

interface RateCardResponse {
  id: UUID;
  provider_id: UUID;
  procedure_code: string;
  procedure_name: string;
  rate_amount: Money;
  effective_date: ISODateTime;
  age_from: number;
  age_to: number;
  gender?: string;
  relationship?: string;
}

interface RegisterProviderRequest {
  name: string;           // required
  type: "hospital" | "clinic" | "pharmacy" | "lab";  // required
  license_number: string; // required
  county?: string;
  address?: string;
  phone: string;          // required
  email: string;          // required, valid email
  contact_person?: string;
}

interface UpdateProviderRequest {
  name?: string;
  county?: string;
  address?: string;
  phone?: string;
  email?: string;
  contact_person?: string;
}

interface CreateContractRequest {
  start_date: ISODateTime;  // required
  end_date: ISODateTime;    // required
  terms?: string;
}

interface CreateRateCardRequest {
  procedure_code: string;   // required
  procedure_name: string;   // required
  rate_amount: Money;        // required, min 1
  effective_date?: ISODateTime;
  age_from?: number;
  age_to?: number;
  gender?: string;
  relationship?: string;
}

interface BulkCreateRateCardRequest {
  rate_cards: CreateRateCardRequest[];  // required, min 1
}

interface UpdateAccreditationRequest {
  accreditation_status: "NONE" | "PENDING" | "ACCREDITED" | "EXPIRED" | "REVOKED";  // required
  accreditation_expiry?: string;  // YYYY-MM-DD
  accreditation_body?: string;
}
```

### Reinsurance Types

```typescript
// lib/types/reinsurance.ts

interface TreatyResponse {
  id: UUID;
  treaty_number: string;
  name: string;
  treaty_type: "QUOTA_SHARE" | "XOL";
  status: "DRAFT" | "ACTIVE" | "EXPIRED" | "TERMINATED";
  effective_date: ISODateTime;
  expiry_date: ISODateTime;
  retention_limit: Money;
  currency: string;
  notes?: string;
  created_by: UUID;
  created_at: ISODateTime;
  updated_at: ISODateTime;
}

interface CreateTreatyRequest {
  name: string;            // required
  treaty_type: "QUOTA_SHARE" | "XOL";  // required
  effective_date: ISODateTime;          // required
  expiry_date: ISODateTime;             // required
  retention_limit?: Money;
  currency?: string;
  notes?: string;
}

interface UpdateTreatyRequest {
  name?: string;
  effective_date?: ISODateTime;
  expiry_date?: ISODateTime;
  retention_limit?: Money;
  currency?: string;
  notes?: string;
}

interface TreatyDetailResponse extends TreatyResponse {
  participants: TreatyParticipantResponse[];
  layers: TreatyLayerResponse[];
  profit_commission_rules: ProfitCommissionResponse[];
}

interface TreatyParticipantResponse {
  id: UUID;
  treaty_id: UUID;
  reinsurer_name: string;
  share_percentage: number;
  commission_rate: number;
  is_lead: boolean;
  created_at: ISODateTime;
  updated_at: ISODateTime;
}

interface TreatyLayerResponse {
  id: UUID;
  treaty_id: UUID;
  layer_number: number;
  attachment_point: Money;
  layer_limit: Money;
  deductible_amount: Money;
  premium_rate: number;
  aggregate_limit?: Money;
  aggregate_used: Money;
  created_at: ISODateTime;
  updated_at: ISODateTime;
}

interface AddParticipantRequest {
  reinsurer_name: string;         // required
  share_percentage: number;       // required, >0 and <=100
  commission_rate?: number;       // >=0
  is_lead?: boolean;
}

interface UpdateParticipantRequest {
  reinsurer_name?: string;
  share_percentage?: number;      // >=0 and <=100
  commission_rate?: number;       // >=0
  is_lead?: boolean;
}

interface AddLayerRequest {
  layer_number: number;           // required, >=1
  attachment_point: Money;        // required, >=0
  layer_limit: Money;             // required, >0
  deductible_amount?: Money;      // >=0
  premium_rate?: number;          // >=0
  aggregate_limit?: Money;
}

interface UpdateLayerRequest {
  attachment_point?: Money;
  layer_limit?: Money;
  deductible_amount?: Money;
  premium_rate?: number;
  aggregate_limit?: Money;
}

interface AddProfitCommissionRuleRequest {
  commission_type: "SLIDING_SCALE" | "FLAT" | "CARRY_FORWARD";  // required
  loss_ratio_from?: number;  // >=0
  loss_ratio_to?: number;    // >=0
  commission_rate: number;   // required, >=0
  carry_forward_years?: number;
}

interface CedePremiumRequest {
  treaty_id: UUID;  // required
  policy_id: UUID;  // required
  amount: Money;    // required, >0
}

interface AutoCedePolicyPremiumRequest {
  policy_id: UUID;  // required
  amount: Money;    // required, >0
}

interface CessionResponse {
  id: UUID;
  cession_number: string;
  treaty_id: UUID;
  policy_id: UUID;
  treaty_layer_id?: UUID;
  cession_type: "PREMIUM" | "CLAIM";
  gross_amount: Money;
  ceded_amount: Money;
  retained_amount: Money;
  commission_amount: Money;
  share_percentage: number;
  status: "PENDING" | "BOOKED" | "REVERSED";
  created_by: UUID;
  created_at: ISODateTime;
  updated_at: ISODateTime;
}

interface CreateRecoveryRequest {
  claim_id: UUID;             // required
  treaty_id: UUID;            // required
  treaty_layer_id?: UUID;
  cession_id?: UUID;
  gross_amount: Money;        // required, >0
  recoverable_amount: Money;  // required, >0
  notes?: string;
}

interface ApplyRecoveryForClaimRequest {
  approved_amount: Money;  // required, >0
}

interface RecoveryWorkflowRequest {
  notes?: string;  // used for acknowledge, request-info, approve transitions
}

interface RecordPaymentRequest {
  amount: Money;  // required, >0
  notes?: string;
}

interface RecoveryResponse {
  id: UUID;
  recovery_number: string;
  claim_id: UUID;
  treaty_id: UUID;
  treaty_layer_id?: UUID;
  cession_id?: UUID;
  gross_claim_amount: Money;
  recoverable_amount: Money;
  recovered_amount: Money;
  outstanding_amount: Money;
  status: "NOTIFIED" | "ACKNOWLEDGED" | "INFO_REQUESTED" | "APPROVED" | "PAID" | "WRITTEN_OFF";
  workflow_status: "NOTIFICATION" | "ACKNOWLEDGMENT" | "INFO_REQUEST" | "APPROVAL" | "PAYMENT";
  notes?: string;
  created_by: UUID;
  created_at: ISODateTime;
  updated_at: ISODateTime;
}

interface RecoveryWorkflowEventResponse {
  id: UUID;
  recovery_id: UUID;
  from_status: string;
  to_status: string;
  event_type: string;
  notes?: string;
  performed_by: UUID;
  created_at: ISODateTime;
}

interface RecoveryDetailResponse extends RecoveryResponse {
  workflow_events: RecoveryWorkflowEventResponse[];
}

interface GenerateBordereauRequest {
  treaty_id: UUID;            // required
  period_start: ISODateTime;  // required
  period_end: ISODateTime;    // required
}

interface BordereauResponse {
  id: UUID;
  bordereau_number: string;
  treaty_id: UUID;
  bordereau_type: "PREMIUM" | "CLAIM";
  period_start: ISODateTime;
  period_end: ISODateTime;
  total_gross: Money;
  total_ceded: Money;
  total_commission: Money;
  item_count: number;
  status: "DRAFT" | "FINALIZED" | "SENT";
  created_by: UUID;
  created_at: ISODateTime;
  updated_at: ISODateTime;
}

interface BordereauItemResponse {
  id: UUID;
  bordereau_id: UUID;
  cession_id?: UUID;
  recovery_id?: UUID;
  policy_number?: string;
  claim_number?: string;
  gross_amount: Money;
  ceded_amount: Money;
  commission_amount: Money;
  created_at: ISODateTime;
}

interface BordereauDetailResponse extends BordereauResponse {
  items: BordereauItemResponse[];
}

interface GenerateStatementRequest {
  treaty_id: UUID;            // required
  participant_id: UUID;       // required
  period_start: ISODateTime;  // required
  period_end: ISODateTime;    // required
}

interface CalculateProfitCommissionRequest {
  treaty_id: UUID;            // required
  period_start: ISODateTime;  // required
  period_end: ISODateTime;    // required
}

interface ReinsurerStatementResponse {
  id: UUID;
  statement_number: string;
  treaty_id: UUID;
  participant_id: UUID;
  period_start: ISODateTime;
  period_end: ISODateTime;
  premium_ceded: Money;
  claims_recovered: Money;
  commission_due: Money;
  profit_commission: Money;
  net_balance: Money;
  status: "DRAFT" | "ISSUED" | "ACKNOWLEDGED" | "SETTLED";
  created_by: UUID;
  created_at: ISODateTime;
  updated_at: ISODateTime;
}

interface ProfitCommissionResponse {
  id: UUID;
  treaty_id: UUID;
  commission_type: "SLIDING_SCALE" | "FLAT" | "CARRY_FORWARD";
  loss_ratio_from: number;
  loss_ratio_to: number;
  commission_rate: number;
  carry_forward_years: number;
  carry_forward_balance: Money;
  period_start?: ISODateTime;
  period_end?: ISODateTime;
  calculated_amount: Money;
  created_at: ISODateTime;
  updated_at: ISODateTime;
}

interface ProfitCommissionCalculationResponse {
  treaty_id: UUID;
  premium_ceded: Money;
  claims_recovered: Money;
  loss_ratio: number;
  net_profit: Money;
  commission_rate: number;
  commission_amount: Money;
  carry_forward: Money;
}

interface TreatyAlertResponse {
  id: UUID;
  treaty_id: UUID;
  treaty_layer_id?: UUID;
  alert_type: "LIMIT_BREACH" | "AGGREGATE_WARNING" | "CATASTROPHE_THRESHOLD" | "EXPIRY_WARNING";
  severity: "LOW" | "MEDIUM" | "HIGH" | "CRITICAL";
  message: string;
  threshold_value: Money;
  current_value: Money;
  is_acknowledged: boolean;
  acknowledged_by?: UUID;
  acknowledged_at?: ISODateTime;
  created_at: ISODateTime;
}

interface AgedRecoveryBucketResponse {
  bucket: string;
  count: number;
  total_outstanding: Money;
}

interface ReinsuranceDashboardResponse {
  active_treaty_count: number;
  total_ceded_premiums: Money;
  total_recoverable: Money;
  total_recovered: Money;
  total_outstanding: Money;
  cession_ratio: number;
  recovery_success_rate: number;
  unacknowledged_alerts: number;
}
```

### Notification & Audit Types

```typescript
// lib/types/notification.ts

interface NotificationResponse {
  id: UUID;
  user_id: UUID;
  channel: "SMS" | "EMAIL" | "IN_APP" | "PUSH";
  type: "QUOTATION" | "APPROVAL" | "CLAIM" | "POLICY" | "DOCUMENT";
  subject: string;
  body: string;
  metadata?: Record<string, unknown>;
  status: "PENDING" | "SENT" | "DELIVERED" | "FAILED" | "READ";
  retry_count: number;
  max_retries: number;
  sent_at?: ISODateTime;
  read_at?: ISODateTime;
  created_at: ISODateTime;
  updated_at: ISODateTime;
}

// lib/types/audit.ts

interface AuditEventResponse {
  id: UUID;
  user_id: UUID;
  entity_type: string;
  entity_id: UUID;
  action: "CREATE" | "UPDATE" | "DELETE" | "STATE_CHANGE";
  old_value?: Record<string, unknown>;
  new_value?: Record<string, unknown>;
  ip_address: string;
  user_agent: string;
  created_at: ISODateTime;
}

// lib/types/analytics.ts

interface DashboardResponse {
  claims_volume: ClaimsVolumeResponse;
  approval_rate: number;
  average_tat_hours: number;
  loss_ratio: number;
  fraud_rate: number;
  total_premium_collected: Money;
  total_claims_paid: Money;
  top_providers: TopProviderResponse[];
}

interface ClaimsVolumeResponse {
  total_claims: number;
  approved_claims: number;
  rejected_claims: number;
  manual_review_claims: number;
  paid_claims: number;
}

interface TopProviderResponse {
  id: UUID;
  name: string;
  claim_count: number;
  total_amount: Money;
  total_approved: Money;
}

interface KPIResponse {
  approval_rate: number;
  average_tat_hours: number;
  loss_ratio: number;
  fraud_rate: number;
  total_premium_collected: Money;
  total_claims_paid: Money;
}
```

---

## 6. API Client Setup

### Axios Instance

```typescript
// lib/api/client.ts

import axios from "axios";

const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL + "/api/v1",
  headers: { "Content-Type": "application/json" },
});

// Request interceptor: inject access token
apiClient.interceptors.request.use((config) => {
  const token = getAccessToken(); // from store/cookies
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor: handle 401 → refresh token
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      try {
        await refreshAccessToken();
        return apiClient(error.config); // retry original request
      } catch {
        redirectToLogin();
      }
    }
    return Promise.reject(error);
  }
);
```

### API Method Pattern

```typescript
// lib/api/claims.ts (example)

export const claimsApi = {
  list: (params: { page: number; page_size: number; status?: string }) =>
    apiClient.get<PaginatedResponse<ClaimResponse>>("/claims", { params }),

  get: (id: string) =>
    apiClient.get<ApiResponse<ClaimResponse>>(`/claims/${id}`),

  submit: (data: SubmitClaimRequest) =>
    apiClient.post<ApiResponse<ClaimResponse>>("/claims", data),

  vet: (id: string, data: { vetted_amount: number; notes?: string }) =>
    apiClient.put<ApiResponse<ClaimResponse>>(`/claims/${id}/vet`, data),

  approve: (id: string, data: { decision: string; reason?: string }) =>
    apiClient.put<ApiResponse<ClaimResponse>>(`/claims/${id}/approve`, data),

  reject: (id: string, data: { decision: string; reason?: string }) =>
    apiClient.put<ApiResponse<ClaimResponse>>(`/claims/${id}/reject`, data),

  markReadyForPayment: (id: string) =>
    apiClient.put<ApiResponse<ClaimResponse>>(`/claims/${id}/ready-for-payment`),

  markPaid: (id: string) =>
    apiClient.put<ApiResponse<ClaimResponse>>(`/claims/${id}/mark-paid`),

  listSLABreached: (params: { page: number; page_size: number }) =>
    apiClient.get<PaginatedResponse<ClaimResponse>>("/claims/sla-breached", { params }),

  bulkSubmit: (data: { claims: SubmitClaimRequest[] }) =>
    apiClient.post<ApiResponse<BulkClaimResultResponse>>("/claims/bulk", data),
};
```

---

## 7. Utility Functions

### Money Formatting

```typescript
// lib/utils/money.ts

export function formatMoney(cents: number, currency = "KES"): string {
  const amount = cents / 100;
  return `${currency} ${amount.toLocaleString("en-KE", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })}`;
}

// formatMoney(800000) → "KES 8,000.00"
// formatMoney(1234567) → "KES 12,345.67"

export function parseMoney(display: string): number {
  const cleaned = display.replace(/[^0-9.]/g, "");
  return Math.round(parseFloat(cleaned) * 100);
}
```

### Date Formatting

```typescript
// lib/utils/date.ts

import { format, formatDistanceToNow } from "date-fns";

export function formatDate(iso: string): string {
  return format(new Date(iso), "dd MMM yyyy");
}

export function formatDateTime(iso: string): string {
  return format(new Date(iso), "dd MMM yyyy HH:mm");
}

export function timeAgo(iso: string): string {
  return formatDistanceToNow(new Date(iso), { addSuffix: true });
}
```

### Status Badge Colors

See **Section 5 → "Status Badge Color Mapping (Complete)"** for the comprehensive `statusColors` object covering all 75+ status values with `getStatusColor()` helper.

---

## 8. Role-Based Access Control

### Page Visibility by Role

| Page / Feature | Admin | Manager | Underwriter | ClaimsOfficer | Finance | SalesAgent | Provider | Member |
|---|---|---|---|---|---|---|---|---|
| Dashboard | Y | Y | Y | Y | Y | Y | Y | Y |
| Products (Plans) | Y | Y | Y | - | - | - | - | - |
| Leads | Y | Y | - | - | - | Y | - | - |
| Quotations | Y | Y | Y | - | - | Y | - | - |
| Approval Limits | Y | - | - | - | - | - | - | - |
| Policies | Y | Y | Y | Y | Y | Y | - | - |
| Members | Y | Y | Y | Y | - | - | - | - |
| Endorsements | Y | Y | Y | - | - | - | - | - |
| Renewals | Y | Y | Y | - | - | - | - | - |
| Underwriting | Y | - | Y | - | - | - | - | - |
| Claims | Y | Y | Y | Y | Y | - | Y | Y |
| Pre-Auth | Y | Y | Y | Y | - | - | Y | - |
| Cases | Y | Y | - | Y | - | - | Y | - |
| Billing (Invoices) | Y | Y | - | - | Y | - | - | - |
| Payments | Y | - | - | - | Y | - | - | - |
| Remittances | Y | - | - | - | Y | - | - | - |
| Reinsurance | Y | Y | - | - | Y | - | - | - |
| Providers | Y | Y | Y | - | - | - | Y | - |
| Users | Y | - | - | - | - | - | - | - |
| Audit Trail | Y | Y | - | - | - | - | - | - |
| Analytics | Y | Y | Y | Y | Y | Y | - | - |
| Notifications | Y | Y | Y | Y | Y | Y | Y | Y |

### Action Permissions by Role

| Action | Roles |
|--------|-------|
| Claim: Vet | Admin, ClaimsOfficer |
| Claim: Approve/Reject | Admin, Manager |
| Claim: Ready for Payment | Admin, Finance |
| Claim: Mark Paid/Part Paid | Admin, Finance |
| Claim: Import CSV | Admin, ClaimsOfficer |
| Quotation: Approve Version | Admin, Underwriter, Manager |
| Quotation: Expire (batch) | Admin |
| Endorsement: Approve/Reject | Admin, Manager |
| Renewal: Expire (batch) | Admin |
| Underwriting: Review | Admin, Underwriter |
| UW Flag: Resolve/Override | Admin, Underwriter |
| Credit Note: Approve/Apply | Admin |
| Approval Limits: CRUD | Admin |
| Users: CRUD | Admin |

### Implementation

```typescript
// lib/utils/constants.ts

export const ROLE_PERMISSIONS: Record<string, string[]> = {
  Admin: ["*"],
  Manager: ["claims.approve", "quotations.approve", "endorsements.approve", "renewals.approve"],
  Underwriter: ["underwriting.review", "quotations.approve", "flags.resolve"],
  ClaimsOfficer: ["claims.vet", "claims.import"],
  Finance: ["claims.pay", "payments.manage", "remittances.manage"],
  SalesAgent: ["leads.manage", "quotations.manage"],
  Provider: ["claims.view", "preauths.view", "cases.view"],
  Member: ["claims.view"],
};

// Helper to check if current user can perform action
export function canPerform(userRole: string, action: string): boolean {
  if (userRole === "Admin") return true;
  const perms = ROLE_PERMISSIONS[userRole] || [];
  return perms.includes(action) || perms.includes("*");
}
```

---

## 9. UI Component Specifications

### Data Table

- Server-side pagination (send `page` and `page_size` to API)
- Sortable column headers
- Filter toolbar with dropdowns and search input
- Row click → navigate to detail page
- Checkbox selection for bulk actions
- Empty state with message
- Loading skeleton while fetching

### Status Badge

```tsx
<StatusBadge status="APPROVED" />
// Renders: <span class="bg-green-100 text-green-800 px-2 py-1 rounded-full text-xs font-medium">APPROVED</span>
```

### Money Display

```tsx
<MoneyDisplay amount={800000} />
// Renders: "KES 8,000.00"
```

### Confirm Dialog

Used for all destructive actions (terminate, reject, remove member, write-off, etc.):
- Title, description, confirm button (red for destructive), cancel button
- Optional reason/notes textarea input

### File Upload

- Drag-and-drop zone
- File type restrictions (CSV for imports, any for documents)
- Progress indicator
- File size display
- Used for: member CSV import, claim CSV import, claim documents, quotation documents, provider statements

### Toast Notifications

- Success: green, auto-dismiss after 5s
- Error: red, requires manual dismiss
- Shown for: form submission, status changes, bulk operations

### Charts

| Chart | Library | Data Source |
|-------|---------|------------|
| Claims Volume (bar) | Recharts BarChart | `GET /analytics/dashboard` |
| Premium Trend (line) | Recharts LineChart | Custom aggregation |
| Loss Ratio (pie) | Recharts PieChart | `GET /analytics/dashboard` |
| Aged Recovery (bar) | Recharts BarChart | `GET /recoveries/aged-analysis` |
| Aggregate Usage (progress) | Custom component | Treaty layer aggregate_used / aggregate_limit |

---

## 10. Key Business Rules for Frontend

### Money Input Fields
- Accept user input in "whole" currency (e.g., user types "8000" for 8000 KES)
- Convert to cents before sending to API: `value × 100`
- Display from API in currency format: `cents / 100`
- All monetary amounts from API are int64 cents (BIGINT)

### SLA Tracking
- Claims have `sla_breach_at` timestamp (created_at + 48 hours)
- Calculate remaining time: `sla_breach_at - now`
- Display countdown (e.g., "23h 45m remaining")
- If breached (now > sla_breach_at): show "BREACHED" in red background
- Approaching (< 24h remaining): show in yellow/amber
- Scheduler checks every 4 hours and sends IN_APP notifications

### Quotation Pricing
- Base premium comes from plan's calculate-premium endpoint (uses age-band matching from proposed_members DOB)
- **Discount/Loading in basis points**: percentage values use `/10000` divisor
  - 500 bps = 5%, 1000 bps = 10%, 5000 bps = 50%
  - Fixed values are in cents
- Final = base - discountAmount + loadingAmount (floored at 0)
- **Approval limits**: If no limits configured for user's role → auto requires approval. If discount/loading exceeds limits → requires_approval = true
- **At approval time**: Approver's limits also checked. If exceeds → error with escalation to next role
- Quotation validity: 30 days from issuance

### Quotation → Policy Conversion
When converting accepted quotation:
1. Creates policy (DRAFT) with client info
2. Enrolls proposed_members from latest version
3. Creates installment schedule with version's billing_frequency
4. Quotation → CONVERTED, Lead → WON

### Premium Display
- Always show in cents converted to currency
- For installments: divide by frequency (monthly=12, quarterly=4, semi_annual=2, annual=1)

### Claim Adjudication — What the Frontend Should Display
The adjudication decision contains detailed rule results. Display these as an expandable panel:
- Each rule has: category (eligibility/coverage/limits/fraud), name, result (PASS/FAIL/FLAG), details
- **Key calculation chain** (show in order):
  1. Annual limit: remaining = limit - used this year
  2. Sub-limit: per_visit or per_item cap
  3. Deductible: subtracted from payable
  4. Co-pay: percentage (payable × rate/100) or fixed amount
  5. Pre-auth cap: payable capped at approved amount
  6. Member responsibility = total - payable

### Fraud Flags Display
Color-code by severity:
- LOW (yellow): REPEAT_VISIT
- MEDIUM (orange): FREQUENCY, RATE_CARD_OVERCHARGE
- HIGH (red): AMOUNT_THRESHOLD, EXPIRED_CONTRACT
- CRITICAL (purple): SUSPENDED_PROVIDER

### Claim Vetting — Type-Specific Rules
Display validation messages in the vet form:
- **DIRECT**: If inpatient (has admission_date), warn if no pre-auth reference
- **REIMBURSEMENT**: Show validation "Vetted amount cannot exceed total claimed amount"
- **EXCEPTION**: Show validation "Cannot exceed 150% of approved amount" (formula: approved × 1.5)

### Renewal Premium Display
Show the claims-experience loading breakdown:
```
Loss Ratio = (total approved claims / premium) × 100
> 100%: +25% loading (high claims)
> 75%:  +15% loading
> 50%:  +10% loading
< 30%:  -5% discount (good claims)
```

### Member Enrollment
- DOB format: YYYY-MM-DD
- Age calculated from DOB: `now.Year - dob.Year - (1 if birthday not yet passed)`
- **Underwriting checks that may reject enrollment:**
  1. Double insurance: national_id exists on another ACTIVE policy → Error + flag
  2. Age vs premium rules: age outside plan's min_age/max_age for relationship → Error + flag
  3. Plan underwriting rules: MAX_AGE/MIN_AGE evaluated per rule → Flag (may or may not block)
- Show underwriting flags in member detail panel

### Member Removal
Show the pro-rata credit note calculation:
```
totalDays = (policy.EndDate - policy.StartDate) in days
remainingDays = (policy.EndDate - now) in days
premiumDiff = oldPremium - newPremium
refundAmount = premiumDiff × remainingDays / totalDays
```
Credit note auto-approved if reason contains "Pro-rata refund"

### Renewal Member Re-validation
During renewal completion, members may be skipped. Display:
- RENEWAL_SKIP flags with reason (age out of range or double insurance)
- Show which members were NOT copied to the renewed policy

### Endorsement Application
When showing endorsement detail, display the changes JSON appropriately:
- ADD_MEMBER: Show member enrollment form data
- REMOVE_MEMBER: Show member ID and reason
- UPDATE_MEMBER: Show member ID and updated fields
- PLAN_CHANGE: Show new plan ID and reason

### Lead Pipeline Rules
- WON and LOST are terminal states — cannot transition back to NEW
- DORMANT can be set from any state
- When quotation created: lead auto-advances to PROPOSAL_SENT (from NEW/CONTACTED/QUALIFIED)

### Reinsurance Key Formulas for Display

**Cession (Quota Share):**
```
ceded = floor(gross × totalSharePercentage / 100)
retained = gross - ceded
commission = floor(ceded × avgCommissionRate / 100)
If retentionLimit > 0 AND retained < retentionLimit: retained = retentionLimit, recalculate ceded
```

**Recovery (XOL) — per layer:**
```
excess = claimAmount - attachmentPoint
exposure = min(excess, layerLimit)
recoverable = exposure - deductible
Check aggregate: remaining = aggregateLimit - aggregateUsed
```

**Profit Commission:**
```
lossRatio = (claimsRecovered × 100) / premiumCeded
netProfit = premiumCeded - claimsRecovered - carryForwardBalance
If netProfit > 0: commission = floor(netProfit × rate / 100)
If netProfit ≤ 0: carryForward = -netProfit (deficit deferred)
```

**Statement Claims Data Scope:**
> `claims_recovered` in statement generation uses the **all-time cumulative total** of recoveries for the treaty, NOT filtered by the statement period. Display a note on the statement generation form: _"Claims recovered reflects all-time total, not just this statement period."_ This matches the backend behavior in `statement_service_impl.go`.

**Treaty Alerts:**
- Aggregate usage ≥ 80%: HIGH warning
- Aggregate usage ≥ 100%: CRITICAL breach
- Total recoverable > 5M KES (500,000,000 cents): CRITICAL catastrophe
- Expiry within 30 days: MEDIUM warning

### CSV Import Format
**Members CSV:**
- Required: `name`, `date_of_birth` (YYYY-MM-DD), `gender`, `relationship`
- Optional: `national_id`, `phone`, `email`, `kra_pin`, `county`, `address`

**Claims CSV:**
- Required: `policy_id`, `member_id`, `provider_id`, `service_date` (YYYY-MM-DD), `procedure_code`, `procedure_name`, `quantity`, `unit_price`
- Optional: `claim_type` (default "DIRECT"), `diagnosis_code` (default "UNSPECIFIED"; use `;` separator for multiple), `notes`, `preauth_id`
- Quantity defaults to 1 if ≤ 0

Show import preview before confirming. Display results: succeeded count, failed count, error details per line.

### Provider Network & Contract Validation
During adjudication, the system checks:
1. Provider must be in plan's provider network → REJECT if not
2. Provider must have ACTIVE contract covering service date (strict: startDate < serviceDate < endDate) → REJECT if not
3. Provider accreditation checked but only FLAGged (not rejected) if not ACCREDITED

### Pre-Auth Validation in Claims
If a claim references a pre-auth:
1. Pre-auth must be APPROVED (reject if not)
2. Pre-auth must not be expired: ValidityEnd > now (reject if expired)
3. Provider must match (reject if mismatch)
4. Procedure codes checked (FLAG if mismatch, not reject)
5. Claim amount vs approved amount (FLAG if exceeds, not reject)
6. After claim processed: pre-auth status → CLAIMED

### Waiting Period Rules
- Waiting period calculated from **member enrollment date** (member.created_at), NOT policy start date
- Each benefit has a waiting_period_days value
- If claim service_date < (member.created_at + waiting_period_days) → REJECT

### Claim Submission Response Handling
**CRITICAL:** Claim submission (`POST /claims`) returns HTTP **201** even when the claim is auto-rejected by the validator or adjudicator. The frontend must check the response `message` field:
- `"Claim submitted and processed"` → normal success
- `"Claim submitted but rejected: ..."` → auto-rejected (still 201)
- `"Claim submitted, adjudication failed"` → adjudication error, claim stuck in VALIDATED
Show appropriate toast: green for processed, amber for rejected, red for adjudication failure.

### Claim Status Gate Reference
Which actions are available per claim status (for button visibility):

| Status | Available Actions |
|--------|-------------------|
| RECEIVED | Reject |
| VALIDATED | Reject |
| ADJUDICATED | Approve, Reject, Vet |
| MANUAL_REVIEW | Approve, Reject |
| APPROVED | Vet |
| VETTED | Ready for Payment |
| PARTIALLY_VETTED | Ready for Payment |
| READY_FOR_PAYMENT | Mark Paid, Mark Part Paid |
| PAID | (none — terminal) |
| PART_PAID | (none — terminal) |
| REJECTED | (none — terminal) |

### Underwriting Assessment Display
Show the auto-decision interpretation based on risk score:

| Risk Score | Decision | Color |
|-----------|----------|-------|
| ≤ 30 | Auto-Approved | Green |
| 31–60 | Referred for Review | Amber |
| > 60 | Auto-Declined | Red |
| Any (with blocker) | Declined (blocking rule) | Red |

Display the individual rule results with their weight contribution:
```
Rule: MAX_AGE (weight: 15) — TRIGGERED: "Member age 72 exceeds max age 65"
Rule: PRE_EXISTING_CONDITION (weight: 20, blocking) — TRIGGERED: "Pre-existing condition flagged: diabetes = yes"
Total Risk Score: 35 → REFER
```

**Underwriting rule types the frontend should handle:**
- MAX_AGE / MIN_AGE — age-based checks
- DOUBLE_INSURANCE — same national_id on another active policy
- PRE_EXISTING_CONDITION — questionnaire value matches "yes"/"true" (case-insensitive)
- BMI_THRESHOLD — questionnaire["bmi"] exceeds threshold (float comparison)
- WAITING_PERIOD — informational flag only (no rejection)

### Policy Creation Defaults
When creating a policy form:
- `start_date`: default to today
- `end_date`: default to start_date + 1 year
- Premium is auto-set from `plan.base_premium` — display as read-only

### Installment Amount Display
Note: installment amounts use integer division. For policies with premiums not evenly divisible:
```
Monthly installment for 100,001 KES premium:
  100001 / 12 = 8333 (×12 = 99,996 — 5 cents difference)
```
Display the per-installment amount from the API response, NOT a client-side calculation.

### Recovery Payment Behavior
**Important UX note:** `RecordPayment` always transitions the recovery to PAID status, even if the payment amount is less than the recoverable amount. There is no PARTIALLY_PAID state.
- Show a warning in the record payment dialog if `amount < outstanding_amount`: "This payment covers {amount} of {outstanding}. The recovery will be marked as PAID regardless."

### Write-Off Availability
Show the "Write Off" action button on recovery records in ALL statuses **except PAID**:
- Available: NOTIFIED, ACKNOWLEDGED, INFO_REQUESTED, APPROVED
- Hidden: PAID, WRITTEN_OFF

### Treaty Activation Validation
On the treaty detail page, show activation requirements:
- Total participant share must be > 0% and ≤ 100%
- At least one participant must be added
- Partial share is allowed (e.g., 60% — insurer retains 40%)
- Display: "Current share: {totalShare}% — {remainingShare}% retained by insurer"

### Pre-Auth Approval Details
When displaying an approved pre-auth:
- `approved_amount` = `estimated_cost` (system copies the estimate exactly at approval time)
- `auth_code` format: `AUTH-{YEAR}-{6-digit}` (auto-generated)
- `validity_end` = approval date + 30 days

### Credit Note Auto-Approval Display
Credit notes with reason containing "Pro-rata refund" (case-sensitive) are auto-approved immediately upon creation. Show these with an "Auto-Approved" badge instead of the manual approval workflow.

### Document Permission Rules
Quotation documents have role-based edit/delete permissions:
- Default when not specified: `["Admin"]`
- Show edit/delete buttons only if the current user's role is in the respective permission array
- Permission arrays are stored as JSON: `can_edit_roles`, `can_delete_roles`

### Scheduler Status for UI
Only the **Claim SLA Enforcement** scheduler is fully active. All other schedulers (billing cycle, payment retry, policy lapse, pre-auth expiry, notification retry, remittance) are currently stubs.

**Manual expiry actions** (not automated — require button clicks in the UI):
- Expire overdue treaties: `POST /treaties/expire`
- Expire quotations past validity: `POST /quotations/expire`
- Expire pending renewals past expires_at: `POST /renewals/expire`
Add "Run Expiry Check" buttons on the respective list pages (Admin only).

**Caution on "Expire Treaties":** The `POST /treaties/expire` endpoint has a known filter issue — it effectively expires **all** fetched active treaties regardless of actual expiry date. Show a confirmation dialog with extra warning: _"This will expire all active treaties. Are you sure?"_

### Endorsement Changes JSON Structure

When building the endorsement creation form, the `changes` field depends on `endorsement_type`:

| Type | JSON Payload |
|---|---|
| `ADD_MEMBER` | `{ "name": "string", "date_of_birth": "YYYY-MM-DD", "gender": "string", "relationship": "string", "national_id": "string", "phone": "string", "email": "string" }` |
| `REMOVE_MEMBER` | `{ "member_id": "uuid", "reason": "string" }` |
| `UPDATE_MEMBER` | `{ "member_id": "uuid", "updates": { "name": "string", "phone": "string", ... } }` |
| `PLAN_CHANGE` | `{ "new_plan_id": "uuid" }` |

Dynamically render the form fields based on the selected endorsement type. Note: `RejectEndorsement` **overwrites** the original reason — display original reason before rejection if needed.

### Renewal Clarifications

- **Premium rules overwrite claims loading:** The claims experience loading (+25%/+15%/+10%/-5%) is applied first, then premium rules are recalculated. If premium rules return a valid result (> 0), it **completely replaces** the loaded premium — not additive. Display the renewal detail to show both "base premium after loading" and "final premium after rules" for transparency.
- **Rejection reason stored in `premium_change_reason`:** On the renewal detail page, show `premium_change_reason` as the rejection reason for REJECTED renewals.
- **Bulk renewal auto-sets date to 30 days:** `POST /renewals/bulk` auto-sets `renewal_date` to 30 days from now — no date picker needed in the bulk renewal form.
- **New policy is DRAFT:** After completing a renewal, remind the user that the new policy needs to be separately activated.

### Per-Family Premium Short-Circuit

When displaying premium rules on the plan configuration page, add a note: if **any** rule has `calculation_type = "per_family"`, per-member rules are **ignored** during calculation. The per_family rule with the best-fit `min_members` (highest value ≤ member count) wins.

### Pre-Auth No Status Guards

The pre-auth approve/deny endpoints have **no status guards** — a pre-auth can be approved or denied from any status. The UI should still only show approve/deny buttons for `SUBMITTED` and `INFO_REQUESTED` statuses to enforce the intended workflow, even though the backend does not enforce it.

### Invoice Generation Rules

When showing invoice details or generating invoices:
- **Amount** = `policy.PremiumAmount` (full annual premium)
- **Due date** = creation date + 30 days
- **Billing period** = creation date to creation date + 1 month
- Display "Invoice generation is manual — the billing cycle scheduler is not yet active."

### Remittance Rules

- **Period** = auto-set to last 1 month (now - 1 month to now) — not user-configurable
- Only **ACTIVE** providers can have remittances created
- Show "No approved claims for remittance" message if no claims available

### Validation Patterns

Apply these regex patterns for client-side validation (matching backend):
```typescript
const KENYAN_PHONE = /^(?:\+254|254|0)?([17]\d{8})$/;
const EMAIL = /^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$/;
const NATIONAL_ID = /^\d{7,8}$/;

// Normalize phone to +254 format before sending to API
function normalizePhone(phone: string): string {
  const match = phone.match(KENYAN_PHONE);
  if (!match) return phone;
  return `+254${match[1]}`;
}
```

### Pagination Defaults

All paginated API endpoints use these defaults:
```typescript
const PAGINATION_DEFAULTS = {
  page: 1,        // min: 1
  page_size: 20,  // min: 1, max: 100
  sort: "created_at",
  order: "desc",
};
```

### Document S3 Patterns & Notification Triggers

Documents that trigger **IN_APP notifications** (show toast "Document ready"):
- Welcome Letter, Policy Schedule, LOU, Decline Letter

Documents that do **NOT** trigger notifications:
- Member Card, Renewal Notice

### LOU Idempotency

`POST /preauths/:id/lou` is **idempotent** — if a LOU already exists for the pre-auth, it returns the existing document (HTTP 200). The UI should handle this gracefully (show "LOU already generated" instead of creating a duplicate).

### Claim SLA Processing Limits

The SLA enforcement task processes up to **100 claims per cycle** per phase. If there are >100 breached claims, some will not get notifications until the next cycle (every 4 hours). The UI should display SLA status from the claim's `sla_breach_at` field independently — don't rely solely on notifications.

### Statement Integer Truncation

`GenerateStatement` casts `SharePercentage` (float64) to int64 before division. A 12.5% share truncates to 12%. Display the participant's actual `share_percentage` alongside the calculated statement amounts, and note that calculated values may differ slightly from expected due to integer arithmetic.

### Co-Pay Payable Amount Can Go Negative

In the adjudication amount calculation, after co-pay subtraction, `payableAmount` is NOT floored at 0 (unlike the deductible step which has a floor). The frontend must guard against displaying negative payable amounts:
```typescript
// Always floor at 0 for display
const displayPayable = Math.max(0, claim.approved_amount);
// But show warning if backend returned negative
if (claim.approved_amount < 0) {
  showWarning("Payable amount is negative — co-pay exceeded available amount");
}
```

### Fraud Checks — Informational Only

Only the **duplicate claim check** (part of adjudication) affects the claim decision (forces MANUAL_REVIEW). The 6 post-adjudication fraud checks (FREQUENCY, AMOUNT_THRESHOLD, EXPIRED_CONTRACT, SUSPENDED_PROVIDER, RATE_CARD_OVERCHARGE, REPEAT_VISIT) only create `FraudFlag` records — they do NOT change the claim status.

In the claim detail view, display fraud flags as **informational alerts** separate from the adjudication decision. Fraud flags are for manual investigation, not automated outcomes.

### Waiting Period & Age — Checked Against ALL Benefits

Adjudication checks waiting period and age against **ALL** plan benefits, not just the matched one. A claim can be rejected because a benefit category irrelevant to the claim has an unmet waiting period or age restriction. The frontend should display which specific benefit triggered the rejection (available in adjudication decision rule details).

### Provider Statement Reconciliation Details

The reconciliation screen should display the two-phase matching result:

| Match Phase | Description |
|---|---|
| Phase 1 | Matched by claim number |
| Phase 2 | Matched by provider + service date + amount (1 KES tolerance) |
| Unmatched | No match found in either phase |

After reconciliation:
- Matched claims are auto-updated: `PAID` (amounts match within 1 KES) or `PART_PAID` (amounts differ)
- Show discrepancy amounts for each line item (negative = provider undercharged, positive = provider overcharged)
- Filter line items by match status (MATCHED / UNMATCHED / DISPUTED)

### Case Management Transition Conditions

Implement status-based button visibility for case management:

| Transition | Button Label | Visible When | Required Fields |
|---|---|---|---|
| SCHEDULED → ADMITTED | "Admit Patient" | status = SCHEDULED | admission_date |
| ADMITTED → IN_TREATMENT | "Start Treatment" | status = ADMITTED | (none) |
| ADMITTED or IN_TREATMENT → DISCHARGED | "Discharge Patient" | status = ADMITTED or IN_TREATMENT | actual_discharge, actual_cost |
| DISCHARGED → CLOSED | "Close Case" | status = DISCHARGED | (none) |

**Note**: Discharge is allowed from both ADMITTED (early discharge) and IN_TREATMENT states. Close case does NOT verify linked claims are in terminal status.

### Authentication — Critical Frontend Implications

1. **Logout is a no-op on the server** — The backend does NOT invalidate tokens. Implement logout entirely client-side:
   ```typescript
   const logout = () => {
     localStorage.removeItem("access_token");
     localStorage.removeItem("refresh_token");
     // Clear all query cache
     queryClient.clear();
     router.push("/login");
   };
   ```

2. **RefreshToken does NOT check user status** — A suspended user can still refresh tokens. The frontend should check `user.status` on token refresh response and force logout if not ACTIVE.

3. **Register defaults to Member role** — The register page creates Member-role users. Admin must use the Users management page to create other roles.

4. **CreateUser (Admin) does NOT hash passwords** — If building a "Create User" form in the admin panel, warn admin users that the password will not be hashed by this endpoint. Consider calling Register instead, or note this limitation clearly.

### Backend Route-to-Role Mapping (Actual RequireRole Enforcement)

The backend only enforces `RequireRole` on a subset of routes. All other routes are accessible to ANY authenticated user. The frontend should implement **UI-level restrictions** beyond what the backend enforces:

**Backend-enforced restrictions (RequireRole middleware):**

| Action | Backend Required Roles |
|---|---|
| Claim: Import CSV | Admin, ClaimsOfficer |
| Claim: Vet | Admin, ClaimsOfficer |
| Claim: Approve/Reject | Admin, Manager |
| Claim: Ready for Payment | Admin, Finance |
| Claim: Mark Paid/Part Paid | Admin, Finance |
| Quotation: Approve/Reject Version | Admin, Underwriter, Manager |
| Quotation: Expire (batch) | Admin |
| Renewal: Expire (batch) | Admin |
| Underwriting: Review | Admin, Underwriter |
| UW Flag: Resolve/Override | Admin, Underwriter |
| Credit Note: Approve/Apply | Admin |
| Approval Limits: All CRUD | Admin |

**IMPORTANT**: `RequirePermission` middleware exists but is **never used** in any route. Only `RequireRole` is enforced. Roles `Provider`, `Member`, and `SalesAgent` are never restricted by any backend route — they can technically access all non-restricted endpoints. The Page Visibility table in Section 8 represents **recommended** UI restrictions.

### Analytics Implementation Notes

- **Period parameter**: Analytics endpoints accept `period` query param: `week` (7d), `month` (30d, default), `quarter` (90d), `year` (365d)
- **Export CSV is a STUB**: `GET /analytics/export` returns an empty CSV with only headers. Hide this button or show "Coming Soon" badge
- **Reinsurance Analytics**: `GET /analytics/reinsurance` uses hardcoded 365-day period. Do NOT show a period selector for this endpoint — it will be ignored

```typescript
// Period selector options for analytics pages
const ANALYTICS_PERIODS = [
  { label: "Last 7 days", value: "week" },
  { label: "Last 30 days", value: "month" },
  { label: "Last 90 days", value: "quarter" },
  { label: "Last 365 days", value: "year" },
] as const;
```

### Notification System — Fire-and-Forget

Notifications are dispatched via fire-and-forget goroutines on the backend. This means:
- A successful API response does NOT guarantee the notification was delivered
- The notification bell should fetch unread count via `GET /notifications/unread-count` polling (every 30-60 seconds)
- Don't show "notification sent" toasts based on the triggering action's success — instead poll for new notifications

### Benefit Hierarchy — Sub-Benefits

Benefits have a parent-child hierarchy via `parent_benefit_id`. Display as a tree:
```
├── Inpatient (annual limit: 5,000,000 KES)
│   ├── Surgery (sub-limit: 2,000,000 KES per visit)
│   ├── Room & Board (sub-limit: 50,000 KES per day)
│   └── ICU (sub-limit: 100,000 KES per day)
├── Outpatient (annual limit: 500,000 KES)
│   ├── Consultation (sub-limit: 5,000 KES per visit)
│   ├── Lab Tests (sub-limit: 20,000 KES per visit)
│   └── Pharmacy (sub-limit: 10,000 KES per visit)
```

**Coverage check ignores procedure codes** — The backend's `CheckCoverage` method accepts a `procedureCode` parameter but ignores it. Coverage is evaluated only by benefit category (INPATIENT vs OUTPATIENT). Do NOT build procedure-code-level coverage lookup in the UI.

### CheckCoverage Benefit Category Determination

How the backend determines benefit category for a claim:
- If `claim.admission_date` is set → **INPATIENT**
- If `claim.admission_date` is null → **OUTPATIENT**
- If no exact category match → falls back to **first active benefit** (`benefits[0]`)

The claim form should make clear that setting an admission date switches the claim to inpatient category, which affects which benefit is used for coverage and limits.

---

## 11. Environment Variables

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_APP_NAME=HIAS Core
NEXT_PUBLIC_DEFAULT_PAGE_SIZE=20
```

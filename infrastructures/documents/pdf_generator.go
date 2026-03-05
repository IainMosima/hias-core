package documents

import (
	"fmt"
	"time"

	policyEntity "github.com/bitbiz/hias-core/domains/policy/entity"
	preauthEntity "github.com/bitbiz/hias-core/domains/preauth/entity"
	productEntity "github.com/bitbiz/hias-core/domains/product/entity"
	"github.com/go-pdf/fpdf"
)

type PDFGenerator interface {
	GenerateWelcomeLetter(policy *policyEntity.Policy, members []*policyEntity.Member, planName string) ([]byte, error)
	GenerateMemberCard(member *policyEntity.Member, policy *policyEntity.Policy, planName string) ([]byte, error)
	GeneratePolicySchedule(policy *policyEntity.Policy, members []*policyEntity.Member, plan *productEntity.Plan, benefits []*productEntity.Benefit) ([]byte, error)
	GenerateRenewalNotice(policy *policyEntity.Policy, renewal *policyEntity.PolicyRenewal) ([]byte, error)
	GenerateEndorsementLetter(policy *policyEntity.Policy, endorsement *policyEntity.Endorsement) ([]byte, error)
	GenerateLOU(preauth *preauthEntity.PreAuthorization, policy *policyEntity.Policy, memberName, providerName, planName string) ([]byte, error)
	GenerateDeclineLetter(policy *policyEntity.Policy, memberName, claimNumber, rejectionReason string) ([]byte, error)
}

type pdfGenerator struct{}

func NewPDFGenerator() PDFGenerator {
	return &pdfGenerator{}
}

func (g *pdfGenerator) addHeader(pdf *fpdf.Fpdf, title string) {
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, "HIAS Insurance")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(190, 6, "Health Insurance Administration System")
	pdf.Ln(12)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, title)
	pdf.Ln(12)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(6)
}

func (g *pdfGenerator) addFooter(pdf *fpdf.Fpdf) {
	pdf.SetY(-30)
	pdf.SetFont("Arial", "I", 8)
	pdf.Cell(190, 5, fmt.Sprintf("Generated on %s", time.Now().Format("02 Jan 2006 15:04")))
	pdf.Ln(4)
	pdf.Cell(190, 5, "This is a system-generated document. For queries, contact support@hias.co.ke")
}

func (g *pdfGenerator) GenerateWelcomeLetter(policy *policyEntity.Policy, members []*policyEntity.Member, planName string) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	g.addHeader(pdf, "Welcome Letter")

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(190, 7, fmt.Sprintf("Date: %s", time.Now().Format("02 January 2006")))
	pdf.Ln(10)

	pdf.Cell(190, 7, fmt.Sprintf("Dear %s,", policy.PolicyholderName))
	pdf.Ln(10)

	pdf.MultiCell(190, 6, "Welcome to HIAS Insurance. We are pleased to confirm your health insurance policy. Below are the details of your coverage.", "", "", false)
	pdf.Ln(6)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 8, "Policy Details")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	details := [][]string{
		{"Policy Number", policy.PolicyNumber},
		{"Plan", planName},
		{"Status", policy.Status},
		{"Start Date", policy.StartDate.Format("02 Jan 2006")},
		{"End Date", policy.EndDate.Format("02 Jan 2006")},
		{"Annual Premium", fmt.Sprintf("KES %s", formatMoney(policy.PremiumAmount))},
	}
	for _, d := range details {
		pdf.CellFormat(60, 7, d[0]+":", "", 0, "", false, 0, "")
		pdf.CellFormat(130, 7, d[1], "", 0, "", false, 0, "")
		pdf.Ln(7)
	}

	if len(members) > 0 {
		pdf.Ln(8)
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(190, 8, "Enrolled Members")
		pdf.Ln(10)

		pdf.SetFont("Arial", "B", 9)
		pdf.CellFormat(60, 7, "Name", "1", 0, "", false, 0, "")
		pdf.CellFormat(40, 7, "Member Number", "1", 0, "", false, 0, "")
		pdf.CellFormat(40, 7, "Relationship", "1", 0, "", false, 0, "")
		pdf.CellFormat(50, 7, "DOB", "1", 0, "", false, 0, "")
		pdf.Ln(7)

		pdf.SetFont("Arial", "", 9)
		for _, m := range members {
			pdf.CellFormat(60, 7, m.Name, "1", 0, "", false, 0, "")
			pdf.CellFormat(40, 7, m.MemberNumber, "1", 0, "", false, 0, "")
			pdf.CellFormat(40, 7, m.Relationship, "1", 0, "", false, 0, "")
			pdf.CellFormat(50, 7, m.DateOfBirth.Format("02 Jan 2006"), "1", 0, "", false, 0, "")
			pdf.Ln(7)
		}
	}

	g.addFooter(pdf)

	var buf []byte
	w := &byteWriter{buf: &buf}
	err := pdf.OutputAndClose(w)
	return buf, err
}

func (g *pdfGenerator) GenerateMemberCard(member *policyEntity.Member, policy *policyEntity.Policy, planName string) ([]byte, error) {
	pdf := fpdf.New("L", "mm", "A6", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(130, 8, "HIAS Insurance")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 8)
	pdf.Cell(130, 5, "Health Insurance Member Card")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(130, 7, member.Name)
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 9)
	info := [][]string{
		{"Member No", member.MemberNumber},
		{"Policy No", policy.PolicyNumber},
		{"Plan", planName},
		{"Relationship", member.Relationship},
		{"Valid From", policy.StartDate.Format("02 Jan 2006")},
		{"Valid To", policy.EndDate.Format("02 Jan 2006")},
	}
	for _, i := range info {
		pdf.CellFormat(35, 6, i[0]+":", "", 0, "", false, 0, "")
		pdf.CellFormat(95, 6, i[1], "", 0, "", false, 0, "")
		pdf.Ln(6)
	}

	var buf []byte
	w := &byteWriter{buf: &buf}
	err := pdf.OutputAndClose(w)
	return buf, err
}

func (g *pdfGenerator) GeneratePolicySchedule(policy *policyEntity.Policy, members []*policyEntity.Member, plan *productEntity.Plan, benefits []*productEntity.Benefit) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	g.addHeader(pdf, "Policy Schedule")

	pdf.SetFont("Arial", "", 10)
	details := [][]string{
		{"Policy Number", policy.PolicyNumber},
		{"Policyholder", policy.PolicyholderName},
		{"Email", policy.PolicyholderEmail},
		{"Phone", policy.PolicyholderPhone},
		{"Plan", plan.Name},
		{"Status", policy.Status},
		{"Period", fmt.Sprintf("%s to %s", policy.StartDate.Format("02 Jan 2006"), policy.EndDate.Format("02 Jan 2006"))},
		{"Annual Premium", fmt.Sprintf("KES %s", formatMoney(policy.PremiumAmount))},
	}
	for _, d := range details {
		pdf.CellFormat(50, 7, d[0]+":", "", 0, "", false, 0, "")
		pdf.CellFormat(140, 7, d[1], "", 0, "", false, 0, "")
		pdf.Ln(7)
	}

	if len(members) > 0 {
		pdf.Ln(8)
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(190, 8, "Members")
		pdf.Ln(10)

		pdf.SetFont("Arial", "B", 9)
		pdf.CellFormat(50, 7, "Name", "1", 0, "", false, 0, "")
		pdf.CellFormat(35, 7, "Member No", "1", 0, "", false, 0, "")
		pdf.CellFormat(30, 7, "Relationship", "1", 0, "", false, 0, "")
		pdf.CellFormat(25, 7, "Gender", "1", 0, "", false, 0, "")
		pdf.CellFormat(30, 7, "DOB", "1", 0, "", false, 0, "")
		pdf.CellFormat(20, 7, "Status", "1", 0, "", false, 0, "")
		pdf.Ln(7)

		pdf.SetFont("Arial", "", 9)
		for _, m := range members {
			pdf.CellFormat(50, 7, m.Name, "1", 0, "", false, 0, "")
			pdf.CellFormat(35, 7, m.MemberNumber, "1", 0, "", false, 0, "")
			pdf.CellFormat(30, 7, m.Relationship, "1", 0, "", false, 0, "")
			pdf.CellFormat(25, 7, m.Gender, "1", 0, "", false, 0, "")
			pdf.CellFormat(30, 7, m.DateOfBirth.Format("02 Jan 2006"), "1", 0, "", false, 0, "")
			pdf.CellFormat(20, 7, m.Status, "1", 0, "", false, 0, "")
			pdf.Ln(7)
		}
	}

	if len(benefits) > 0 {
		pdf.Ln(8)
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(190, 8, "Benefits Summary")
		pdf.Ln(10)

		pdf.SetFont("Arial", "B", 9)
		pdf.CellFormat(60, 7, "Benefit", "1", 0, "", false, 0, "")
		pdf.CellFormat(35, 7, "Category", "1", 0, "", false, 0, "")
		pdf.CellFormat(40, 7, "Cover Limit", "1", 0, "", false, 0, "")
		pdf.CellFormat(55, 7, "Waiting Period", "1", 0, "", false, 0, "")
		pdf.Ln(7)

		pdf.SetFont("Arial", "", 9)
		for _, b := range benefits {
			pdf.CellFormat(60, 7, b.Name, "1", 0, "", false, 0, "")
			pdf.CellFormat(35, 7, b.Category, "1", 0, "", false, 0, "")
			pdf.CellFormat(40, 7, fmt.Sprintf("KES %s", formatMoney(b.AnnualLimit)), "1", 0, "", false, 0, "")
			pdf.CellFormat(55, 7, fmt.Sprintf("%d days", b.WaitingPeriodDays), "1", 0, "", false, 0, "")
			pdf.Ln(7)
		}
	}

	g.addFooter(pdf)

	var buf []byte
	w := &byteWriter{buf: &buf}
	err := pdf.OutputAndClose(w)
	return buf, err
}

func (g *pdfGenerator) GenerateRenewalNotice(policy *policyEntity.Policy, renewal *policyEntity.PolicyRenewal) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	g.addHeader(pdf, "Policy Renewal Notice")

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(190, 7, fmt.Sprintf("Date: %s", time.Now().Format("02 January 2006")))
	pdf.Ln(10)

	pdf.Cell(190, 7, fmt.Sprintf("Dear %s,", policy.PolicyholderName))
	pdf.Ln(10)

	pdf.MultiCell(190, 6, fmt.Sprintf("Your health insurance policy %s is due for renewal. Please review the details below.", policy.PolicyNumber), "", "", false)
	pdf.Ln(6)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 8, "Current Policy")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(60, 7, "Policy Number:", "", 0, "", false, 0, "")
	pdf.CellFormat(130, 7, policy.PolicyNumber, "", 0, "", false, 0, "")
	pdf.Ln(7)
	pdf.CellFormat(60, 7, "Current Premium:", "", 0, "", false, 0, "")
	pdf.CellFormat(130, 7, fmt.Sprintf("KES %s", formatMoney(policy.PremiumAmount)), "", 0, "", false, 0, "")
	pdf.Ln(7)
	pdf.CellFormat(60, 7, "End Date:", "", 0, "", false, 0, "")
	pdf.CellFormat(130, 7, policy.EndDate.Format("02 Jan 2006"), "", 0, "", false, 0, "")
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 8, "Renewal Details")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(60, 7, "Renewal Date:", "", 0, "", false, 0, "")
	pdf.CellFormat(130, 7, renewal.RenewalDate.Format("02 Jan 2006"), "", 0, "", false, 0, "")
	pdf.Ln(7)
	pdf.CellFormat(60, 7, "New Premium:", "", 0, "", false, 0, "")
	pdf.CellFormat(130, 7, fmt.Sprintf("KES %s", formatMoney(renewal.NewPremium)), "", 0, "", false, 0, "")
	pdf.Ln(7)

	if renewal.ExpiresAt != nil {
		pdf.CellFormat(60, 7, "Accept By:", "", 0, "", false, 0, "")
		pdf.CellFormat(130, 7, renewal.ExpiresAt.Format("02 Jan 2006"), "", 0, "", false, 0, "")
		pdf.Ln(7)
	}

	g.addFooter(pdf)

	var buf []byte
	w := &byteWriter{buf: &buf}
	err := pdf.OutputAndClose(w)
	return buf, err
}

func (g *pdfGenerator) GenerateEndorsementLetter(policy *policyEntity.Policy, endorsement *policyEntity.Endorsement) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	g.addHeader(pdf, "Endorsement Confirmation")

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(190, 7, fmt.Sprintf("Date: %s", time.Now().Format("02 January 2006")))
	pdf.Ln(10)

	pdf.Cell(190, 7, fmt.Sprintf("Dear %s,", policy.PolicyholderName))
	pdf.Ln(10)

	pdf.MultiCell(190, 6, fmt.Sprintf("This letter confirms the following endorsement to your policy %s.", policy.PolicyNumber), "", "", false)
	pdf.Ln(6)

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(60, 7, "Endorsement Type:", "", 0, "", false, 0, "")
	pdf.CellFormat(130, 7, endorsement.EndorsementType, "", 0, "", false, 0, "")
	pdf.Ln(7)
	pdf.CellFormat(60, 7, "Effective Date:", "", 0, "", false, 0, "")
	pdf.CellFormat(130, 7, endorsement.EffectiveDate.Format("02 Jan 2006"), "", 0, "", false, 0, "")
	pdf.Ln(7)
	if endorsement.Reason != "" {
		pdf.CellFormat(60, 7, "Reason:", "", 0, "", false, 0, "")
		pdf.CellFormat(130, 7, endorsement.Reason, "", 0, "", false, 0, "")
		pdf.Ln(7)
	}
	if endorsement.PremiumAdjustment != 0 {
		label := "Premium Increase:"
		if endorsement.PremiumAdjustment < 0 {
			label = "Premium Decrease:"
		}
		pdf.CellFormat(60, 7, label, "", 0, "", false, 0, "")
		pdf.CellFormat(130, 7, fmt.Sprintf("KES %s", formatMoney(abs(endorsement.PremiumAdjustment))), "", 0, "", false, 0, "")
		pdf.Ln(7)
	}

	g.addFooter(pdf)

	var buf []byte
	w := &byteWriter{buf: &buf}
	err := pdf.OutputAndClose(w)
	return buf, err
}

func (g *pdfGenerator) GenerateDeclineLetter(policy *policyEntity.Policy, memberName, claimNumber, rejectionReason string) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	g.addHeader(pdf, "Claim Decline Notification")

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(190, 7, fmt.Sprintf("Date: %s", time.Now().Format("02 January 2006")))
	pdf.Ln(10)
	pdf.Cell(190, 7, fmt.Sprintf("Dear %s,", memberName))
	pdf.Ln(10)
	pdf.MultiCell(190, 6, fmt.Sprintf("We regret to inform you that your claim %s under policy %s has been declined.", claimNumber, policy.PolicyNumber), "", "", false)
	pdf.Ln(6)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 8, "Reason for Decline")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(190, 6, rejectionReason, "", "", false)
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(190, 6, "If you believe this decision is in error, you may submit an appeal within 30 days of this notice. Please contact our claims department for further assistance.", "", "", false)

	g.addFooter(pdf)
	var buf []byte
	w := &byteWriter{buf: &buf}
	err := pdf.OutputAndClose(w)
	return buf, err
}

// byteWriter implements io.WriteCloser for fpdf output
type byteWriter struct {
	buf *[]byte
}

func (w *byteWriter) Write(p []byte) (n int, err error) {
	*w.buf = append(*w.buf, p...)
	return len(p), nil
}

func (w *byteWriter) Close() error {
	return nil
}

func formatMoney(cents int64) string {
	whole := cents / 100
	frac := cents % 100
	if frac < 0 {
		frac = -frac
	}
	return fmt.Sprintf("%d.%02d", whole, frac)
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

func (g *pdfGenerator) GenerateLOU(preauth *preauthEntity.PreAuthorization, policy *policyEntity.Policy, memberName, providerName, planName string) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	g.addHeader(pdf, "Letter of Undertaking (LOU)")

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(190, 7, fmt.Sprintf("Date: %s", time.Now().Format("02 January 2006")))
	pdf.Ln(10)

	pdf.Cell(190, 7, fmt.Sprintf("To: %s", providerName))
	pdf.Ln(10)

	pdf.MultiCell(190, 6, fmt.Sprintf("This letter serves as an undertaking by HIAS Insurance to cover the approved medical expenses for the following pre-authorized treatment under policy %s.", policy.PolicyNumber), "", "", false)
	pdf.Ln(6)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 8, "Authorization Details")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	details := [][]string{
		{"Authorization Code", preauth.AuthCode},
		{"Policy Number", policy.PolicyNumber},
		{"Plan", planName},
		{"Member Name", memberName},
		{"Provider", providerName},
		{"Approved Amount", fmt.Sprintf("KES %s", formatMoney(preauth.ApprovedAmount))},
	}
	if preauth.ValidityStart != nil {
		details = append(details, []string{"Valid From", preauth.ValidityStart.Format("02 Jan 2006")})
	}
	if preauth.ValidityEnd != nil {
		details = append(details, []string{"Valid Until", preauth.ValidityEnd.Format("02 Jan 2006")})
	}

	for _, d := range details {
		pdf.CellFormat(60, 7, d[0]+":", "", 0, "", false, 0, "")
		pdf.CellFormat(130, 7, d[1], "", 0, "", false, 0, "")
		pdf.Ln(7)
	}

	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 8, "Terms and Conditions")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	terms := []string{
		"1. This undertaking is valid only for the approved procedures and amounts listed above.",
		"2. Any costs exceeding the approved amount must be pre-approved by HIAS Insurance.",
		"3. The provider must submit claims within 30 days of service delivery.",
		"4. This LOU is non-transferable and applies only to the named member and provider.",
		"5. HIAS Insurance reserves the right to audit and verify all claims submitted under this LOU.",
	}
	for _, t := range terms {
		pdf.MultiCell(190, 6, t, "", "", false)
		pdf.Ln(2)
	}

	g.addFooter(pdf)

	var buf []byte
	w := &byteWriter{buf: &buf}
	err := pdf.OutputAndClose(w)
	return buf, err
}

package scheduler

type SchedulerConfig struct {
	BillingCycleCron       string // Default: "0 0 1 * *" (1st of month)
	PaymentReminderCron    string // Default: "0 8 * * *" (daily 8am)
	PolicyLapseCron        string // Default: "0 1 * * *" (daily 1am)
	PreAuthExpiryCron      string // Default: "0 2 * * *" (daily 2am)
	RemittanceCycleCron    string // Default: "0 0 * * 1" (Monday midnight)
	PaymentRetryCron       string // Default: "0 */4 * * *" (every 4h)
	ReconciliationCron     string // Default: "0 2 * * *" (daily 2am)
	NotificationRetryCron  string // Default: "*/30 * * * *" (every 30m)
	ReportDistributionCron string // Default: "*/5 * * * *" (every 5 minutes)
	ReportCleanupCron      string // Default: "0 2 * * *" (daily 2am)
}

func DefaultSchedulerConfig() SchedulerConfig {
	return SchedulerConfig{
		BillingCycleCron:       "0 0 1 * *",
		PaymentReminderCron:    "0 8 * * *",
		PolicyLapseCron:        "0 1 * * *",
		PreAuthExpiryCron:      "0 2 * * *",
		RemittanceCycleCron:    "0 0 * * 1",
		PaymentRetryCron:       "0 */4 * * *",
		ReconciliationCron:     "0 2 * * *",
		NotificationRetryCron:  "*/30 * * * *",
		ReportDistributionCron: "*/5 * * * *",
		ReportCleanupCron:      "0 2 * * *",
	}
}

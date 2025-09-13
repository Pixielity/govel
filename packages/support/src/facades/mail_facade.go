package facades

import (
	mailInterfaces "govel/types/src/interfaces/mail"
	facade "govel/support/src"
)

// Mail provides a clean, static-like interface to the application's email service.
//
// This facade implements the facade pattern, providing global access to the mail
// service configured in the dependency injection container. It offers a Laravel-style
// API for email operations with automatic service resolution, template rendering,
// attachment handling, queue integration, and multiple mail driver support.
//
// Architecture:
//   - Uses facade.Resolve() internally for service resolution
//   - Automatically caches the resolved mail service for performance
//   - Provides compile-time type safety through generics
//   - Thread-safe for concurrent email operations across goroutines
//   - Supports multiple mail drivers (SMTP, SES, Mailgun, SendGrid, etc.)
//   - Built-in template rendering and email composition
//
// Behavior:
//   - First call: Resolves mail service from container, performs type assertion, caches result
//   - Subsequent calls: Returns cached service instance (extremely fast)
//   - Panics if mail service cannot be resolved (fail-fast behavior)
//   - Automatically handles email formatting, encoding, and delivery
//
// Returns:
//   - MailInterface: The application's mail service instance
//
// Panics:
//   - If no container is set via facades.SetContainer() or support.SetContainer()
//   - If "mail" service is not registered in the container
//   - If the resolved service doesn't implement MailInterface
//   - If container resolution fails for any reason
//
// Performance Characteristics:
//   - First call: ~100-1000ns (depending on container and service complexity)
//   - Subsequent calls: ~10-50ns (cached lookup with atomic operations)
//   - Memory: Minimal overhead, shared cache across all facade calls
//   - Concurrency: Optimized read-write locks minimize contention
//
// Thread Safety:
// This facade is completely thread-safe:
//   - Multiple goroutines can call Mail() concurrently without synchronization
//   - Internal caching uses optimized read-write mutexes
//   - Service resolution is protected against race conditions
//   - Email composition and sending are thread-safe
//
// Usage Examples:
//
//	// Simple text email
//	err := facades.Mail().Send(
//	    "user@example.com",
//	    "Welcome to our service!",
//	    "Thank you for joining us. We're excited to have you!",
//	)
//	if err != nil {
//	    log.Printf("Failed to send email: %v", err)
//	}
//
//	// HTML email with template
//	data := map[string]interface{}{
//	    "name":     "John Doe",
//	    "company":  "Acme Corp",
//	    "loginUrl": "https://app.example.com/login",
//	}
//
//	err := facades.Mail().SendTemplate(
//	    "user@example.com",
//	    "Welcome to {{company}}!",
//	    "emails/welcome",
//	    data,
//	)
//
//	// Email with multiple recipients
//	recipients := []string{
//	    "user1@example.com",
//	    "user2@example.com",
//	    "user3@example.com",
//	}
//
//	err := facades.Mail().SendToMany(
//	    recipients,
//	    "System Maintenance Notice",
//	    "The system will be under maintenance on Sunday at 2 AM UTC.",
//	)
//
//	// Email with CC and BCC
//	mailData := facades.Mail().NewMessage().
//	    To("recipient@example.com").
//	    Cc("manager@example.com").
//	    Bcc("admin@example.com").
//	    Subject("Project Update").
//	    Body("Here's the latest project status...")
//
//	err := facades.Mail().SendMessage(mailData)
//
//	// Email with file attachments
//	message := facades.Mail().NewMessage().
//	    To("client@example.com").
//	    Subject("Invoice #12345").
//	    Body("Please find your invoice attached.").
//	    Attach("/path/to/invoice.pdf").
//	    Attach("/path/to/receipt.jpg")
//
//	err := facades.Mail().SendMessage(message)
//
//	// Email with inline attachments (for HTML emails)
//	htmlBody := `
//	    <h1>Company Newsletter</h1>
//	    <img src="cid:logo" alt="Company Logo">
//	    <p>Welcome to our monthly newsletter!</p>
//	`
//
//	message := facades.Mail().NewMessage().
//	    To("subscriber@example.com").
//	    Subject("Monthly Newsletter").
//	    Html(htmlBody).
//	    Embed("/path/to/logo.png", "logo")
//
//	err := facades.Mail().SendMessage(message)
//
//	// Queued email (background sending)
//	err := facades.Mail().Queue(
//	    "user@example.com",
//	    "Password Reset Request",
//	    "emails/password-reset",
//	    map[string]interface{}{
//	        "resetUrl": "https://app.example.com/reset/abc123",
//	        "expires":  time.Now().Add(24 * time.Hour),
//	    },
//	)
//
//	// Delayed email (send later)
//	sendAt := time.Now().Add(1 * time.Hour)
//	err := facades.Mail().Later(
//	    sendAt,
//	    "user@example.com",
//	    "Reminder: Meeting in 1 hour",
//	    "Don't forget about our scheduled meeting at 3 PM.",
//	)
//
//	// Email with custom headers
//	message := facades.Mail().NewMessage().
//	    To("support@example.com").
//	    Subject("Bug Report").
//	    Body("I found a bug in the application.").
//	    Header("X-Priority", "1").
//	    Header("X-Report-Type", "Bug")
//
//	err := facades.Mail().SendMessage(message)
//
//	// Bulk email with personalization
//	users := []struct {
//	    Email string
//	    Name  string
//	}{
//	    {"user1@example.com", "Alice"},
//	    {"user2@example.com", "Bob"},
//	    {"user3@example.com", "Charlie"},
//	}
//
//	for _, user := range users {
//	    err := facades.Mail().Queue(
//	        user.Email,
//	        "Personal Invitation",
//	        "emails/invitation",
//	        map[string]interface{}{
//	            "name": user.Name,
//	        },
//	    )
//	    if err != nil {
//	        log.Printf("Failed to queue email for %s: %v", user.Email, err)
//	    }
//	}
//
// Advanced Email Patterns:
//
//	// Email service wrapper
//	type EmailService struct {
//	    fromAddress string
//	    fromName    string
//	}
//
//	func NewEmailService(fromAddress, fromName string) *EmailService {
//	    return &EmailService{
//	        fromAddress: fromAddress,
//	        fromName:    fromName,
//	    }
//	}
//
//	func (s *EmailService) SendWelcomeEmail(userEmail, userName string) error {
//	    return facades.Mail().SendTemplate(
//	        userEmail,
//	        "Welcome to our platform, "+userName+"!",
//	        "emails/welcome",
//	        map[string]interface{}{
//	            "name": userName,
//	            "from": s.fromName,
//	        },
//	    )
//	}
//
//	// Email verification workflow
//	func SendVerificationEmail(user User) error {
//	    token := generateVerificationToken(user.ID)
//	    verifyUrl := fmt.Sprintf("https://app.example.com/verify/%s", token)
//
//	    return facades.Mail().SendTemplate(
//	        user.Email,
//	        "Please verify your email address",
//	        "emails/verify-email",
//	        map[string]interface{}{
//	            "name":      user.Name,
//	            "verifyUrl": verifyUrl,
//	            "expires":   time.Now().Add(24 * time.Hour),
//	        },
//	    )
//	}
//
//	// Multi-language email support
//	func SendLocalizedEmail(userEmail, locale, templateKey string, data map[string]interface{}) error {
//	    // Get localized subject
//	    subject := facades.Lang().ForLocale(locale, func() string {
//	        return facades.Lang().Get("email."+templateKey+".subject", data)
//	    })
//
//	    // Determine template path based on locale
//	    templatePath := fmt.Sprintf("emails/%s/%s", locale, templateKey)
//
//	    return facades.Mail().SendTemplate(userEmail, subject, templatePath, data)
//	}
//
//	// Email with retry logic
//	func SendEmailWithRetry(to, subject, body string, maxRetries int) error {
//	    var lastErr error
//
//	    for attempt := 0; attempt <= maxRetries; attempt++ {
//	        err := facades.Mail().Send(to, subject, body)
//	        if err == nil {
//	            return nil // Success
//	        }
//
//	        lastErr = err
//	        if attempt < maxRetries {
//	            // Wait before retry with exponential backoff
//	            time.Sleep(time.Duration(math.Pow(2, float64(attempt))) * time.Second)
//	        }
//	    }
//
//	    return fmt.Errorf("failed to send email after %d attempts: %w", maxRetries+1, lastErr)
//	}
//
//	// Email batch processing
//	func ProcessEmailBatch(emails []EmailJob) {
//	    batchSize := 10
//
//	    for i := 0; i < len(emails); i += batchSize {
//	        end := i + batchSize
//	        if end > len(emails) {
//	            end = len(emails)
//	        }
//
//	        batch := emails[i:end]
//
//	        // Process batch concurrently
//	        var wg sync.WaitGroup
//	        for _, email := range batch {
//	            wg.Add(1)
//	            go func(e EmailJob) {
//	                defer wg.Done()
//	                err := facades.Mail().SendTemplate(e.To, e.Subject, e.Template, e.Data)
//	                if err != nil {
//	                    log.Printf("Failed to send email to %s: %v", e.To, err)
//	                }
//	            }(email)
//	        }
//
//	        wg.Wait()
//
//	        // Rate limiting between batches
//	        time.Sleep(1 * time.Second)
//	    }
//	}
//
//	// Email tracking and analytics
//	type EmailTracker struct {
//	    messageID string
//	    recipient string
//	    sentAt    time.Time
//	}
//
//	func SendTrackedEmail(to, subject, body string) (*EmailTracker, error) {
//	    messageID := generateMessageID()
//
//	    message := facades.Mail().NewMessage().
//	        To(to).
//	        Subject(subject).
//	        Body(body).
//	        Header("Message-ID", messageID).
//	        Header("X-Track-Opens", "true")
//
//	    err := facades.Mail().SendMessage(message)
//	    if err != nil {
//	        return nil, err
//	    }
//
//	    return &EmailTracker{
//	        messageID: messageID,
//	        recipient: to,
//	        sentAt:    time.Now(),
//	    }, nil
//	}
//
//	// Newsletter management
//	func SendNewsletter(newsletter Newsletter, subscribers []Subscriber) error {
//	    for _, subscriber := range subscribers {
//	        // Personalize content
//	        data := map[string]interface{}{
//	            "name":          subscriber.Name,
//	            "preferences":   subscriber.Preferences,
//	            "unsubscribeUrl": generateUnsubscribeUrl(subscriber.ID),
//	        }
//
//	        // Queue email to avoid blocking
//	        err := facades.Mail().Queue(
//	            subscriber.Email,
//	            newsletter.Subject,
//	            "emails/newsletter",
//	            data,
//	        )
//
//	        if err != nil {
//	            log.Printf("Failed to queue newsletter for %s: %v", subscriber.Email, err)
//	        }
//	    }
//
//	    return nil
//	}
//
// Best Practices:
//   - Always use templates for HTML emails to ensure consistent formatting
//   - Queue emails for better performance and reliability
//   - Implement proper error handling and logging
//   - Use meaningful subjects and preview text
//   - Include plain text alternatives for HTML emails
//   - Add unsubscribe links for marketing emails
//   - Validate email addresses before sending
//   - Implement rate limiting to respect provider limits
//
// Email Security:
//   - Use DKIM and SPF records for authentication
//   - Implement DMARC policies
//   - Validate and sanitize all email content
//   - Use TLS for SMTP connections
//   - Implement bounce handling
//   - Monitor for spam complaints
//   - Use reputation monitoring services
//
// Error Handling:
// This facade uses panic-on-error behavior for clean code:
//   - Most application code can assume mail service always works
//   - Failures are detected early and halt execution
//   - No need for error checking in normal application flow
//   - Container configuration issues are caught immediately
//
// Alternative Error-Safe Access:
// If you need error handling instead of panics, use support package directly:
//
//	mail, err := facade.TryResolve[MailInterface]("mail")
//	if err != nil {
//	    // Handle mail service unavailability gracefully
//	    log.Printf("Mail service unavailable: %v", err)
//	    return // Skip sending email
//	}
//	err = mail.Send(to, subject, body)
//
// Testing Support:
// This facade supports comprehensive testing through service swapping:
//
//	func TestEmailSending(t *testing.T) {
//	    // Create a test mail service that captures sent emails
//	    testMail := &TestMailService{
//	        sentEmails: []SentEmail{},
//	    }
//
//	    // Swap the real mail service with test service
//	    restore := support.SwapService("mail", testMail)
//	    defer restore() // Always restore after test
//
//	    // Now facades.Mail() returns testMail
//	    emailService := NewEmailService("test@example.com", "Test App")
//
//	    err := emailService.SendWelcomeEmail("user@example.com", "John Doe")
//	    require.NoError(t, err)
//
//	    // Verify email was "sent"
//	    emails := testMail.GetSentEmails()
//	    assert.Len(t, emails, 1)
//	    assert.Equal(t, "user@example.com", emails[0].To)
//	    assert.Contains(t, emails[0].Subject, "Welcome")
//	}
//
// Container Configuration:
// Ensure the mail service is properly configured in your container:
//
//	// Example mail registration
//	container.Singleton("mail", func() interface{} {
//	    config := mail.Config{
//	        // Default driver
//	        Driver: "smtp", // smtp, ses, mailgun, sendgrid, log
//
//	        // From address configuration
//	        From: mail.Address{
//	            Address: "noreply@example.com",
//	            Name:    "Example App",
//	        },
//
//	        // SMTP configuration
//	        SMTP: mail.SMTPConfig{
//	            Host:     "smtp.gmail.com",
//	            Port:     587,
//	            Username: "your-email@gmail.com",
//	            Password: "your-app-password",
//	            Encryption: "tls", // tls, ssl, none
//	        },
//
//	        // AWS SES configuration
//	        SES: mail.SESConfig{
//	            Region:    "us-east-1",
//	            AccessKey: "your-access-key",
//	            SecretKey: "your-secret-key",
//	        },
//
//	        // Mailgun configuration
//	        Mailgun: mail.MailgunConfig{
//	            Domain:    "mg.example.com",
//	            APIKey:    "your-api-key",
//	            Endpoint:  "https://api.mailgun.net/v3",
//	        },
//
//	        // Template configuration
//	        Templates: mail.TemplateConfig{
//	            Path:   "./templates/emails",
//	            Engine: "html/template", // or "pongo2", "handlebars"
//	            Cache:  true,
//	        },
//
//	        // Queue configuration
//	        Queue: mail.QueueConfig{
//	            Connection: "default",
//	            Queue:      "emails",
//	            Retry:      3,
//	            Timeout:    time.Minute * 5,
//	        },
//
//	        // Logging configuration
//	        Log: mail.LogConfig{
//	            Channel: "mail",
//	            Level:   "info",
//	        },
//	    }
//
//	    mailService, err := mail.NewMailService(config)
//	    if err != nil {
//	        log.Fatalf("Failed to create mail service: %v", err)
//	    }
//
//	    return mailService
//	})
func Mail() mailInterfaces.MailInterface {
	// Use facade.Resolve() for clean facade implementation:
	// - Resolves "mail" service from the dependency injection container
	// - Performs type assertion to MailInterface
	// - Caches the result for subsequent calls
	// - Panics with descriptive error if resolution fails
	// - Thread-safe with optimized locking
	return facade.Resolve[mailInterfaces.MailInterface](mailInterfaces.MAIL_TOKEN)
}

// MailWithError provides error-safe access to the mail service.
//
// This function offers the same functionality as Mail() but returns errors
// instead of panicking, making it suitable for error-sensitive contexts where
// you want to handle mail service unavailability gracefully.
//
// This is a convenience wrapper around facade.TryResolve() that provides
// the same caching and performance benefits as Mail() but with error handling.
//
// Returns:
//   - MailInterface: The resolved mail instance (nil if error occurs)
//   - error: Detailed error information if resolution fails
//
// Errors:
//   - support.FacadeError: If container not set or service resolution fails
//   - Type assertion errors: If service doesn't implement MailInterface
//
// Usage Examples:
//
//	// Basic error-safe email sending
//	mail, err := facades.MailWithError()
//	if err != nil {
//	    log.Printf("Mail service unavailable: %v", err)
//	    // Maybe log the email content for later manual processing
//	    return fmt.Errorf("unable to send email")
//	}
//	err = mail.Send("user@example.com", "Welcome!", "Thank you for joining us.")
//
//	// Conditional email sending
//	if mail, err := facades.MailWithError(); err == nil {
//	    // Send optional notification emails
//	    mail.Queue("admin@example.com", "System Alert", alertMessage)
//	}
//
//	// Health check pattern
//	func CheckMailHealth() error {
//	    mail, err := facades.MailWithError()
//	    if err != nil {
//	        return fmt.Errorf("mail service unavailable: %w", err)
//	    }
//
//	    // Test basic mail functionality
//	    err = mail.Send("test@example.com", "Health Check", "Test message")
//	    if err != nil {
//	        return fmt.Errorf("mail service not working properly: %w", err)
//	    }
//
//	    return nil
//	}
func MailWithError() (mailInterfaces.MailInterface, error) {
	// Use facade.TryResolve() for error-return behavior:
	// - Resolves "mail" service from the dependency injection container
	// - Performs type assertion with error handling
	// - Caches the result for subsequent calls
	// - Returns detailed error information instead of panicking
	// - Thread-safe with optimized locking
	return facade.TryResolve[mailInterfaces.MailInterface](mailInterfaces.MAIL_TOKEN)
}

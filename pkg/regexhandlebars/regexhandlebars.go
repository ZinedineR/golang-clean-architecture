package regexhandlebars

import (
	"boiler-plate-clean/internal/entity"
	"regexp"
	"sort"
)

func RemoveDuplicatesInPlace(slice []string) []string {
	// if there are 0 or 1 items we return the slice itself.
	if len(slice) < 2 {
		return slice
	}

	// make the slice ascending sorted.
	sort.SliceStable(slice, func(i, j int) bool { return slice[i] < slice[j] })

	uniqPointer := 0

	for i := 1; i < len(slice); i++ {
		// compare a current item with the item under the unique pointer.
		// if they are not the same, write the item next to the right of the unique pointer.
		if slice[uniqPointer] != slice[i] {
			uniqPointer++
			slice[uniqPointer] = slice[i]
		}
	}

	return slice[:uniqPointer+1]
}

func ExtractMailTemplate(template entity.MailTemplate) []string {
	placeholders := make([]string, 0)

	// Regular expression to find Handlebars placeholders
	re := regexp.MustCompile(`{{[^{}]+}}`)

	// Extract placeholders from Subject
	matches := re.FindAllString(template.Subject, -1)
	for _, match := range matches {
		placeholders = append(placeholders, match)
	}

	// Extract placeholders from Body
	matches = re.FindAllString(template.Body, -1)
	for _, match := range matches {
		placeholders = append(placeholders, match)
	}
	return RemoveDuplicatesInPlace(placeholders)
}

func ExtractSMSTemplate(template entity.SMSTemplate) []string {
	placeholders := make([]string, 0)

	// Regular expression to find Handlebars placeholders
	re := regexp.MustCompile(`{{[^{}]+}}`)

	// Extract placeholders from Body
	matches := re.FindAllString(template.Body, -1)
	for _, match := range matches {
		placeholders = append(placeholders, match)
	}

	return RemoveDuplicatesInPlace(placeholders)
}

func ExtractNotificationTemplate(template entity.NotificationTemplate) []string {
	placeholders := make([]string, 0)

	// Regular expression to find Handlebars placeholders
	re := regexp.MustCompile(`{{[^{}]+}}`)
	// Extract placeholders from Subject
	matches := re.FindAllString(template.Subject, -1)
	for _, match := range matches {
		placeholders = append(placeholders, match)
	}
	// Extract placeholders from Body
	matches = re.FindAllString(template.Body, -1)
	for _, match := range matches {
		placeholders = append(placeholders, match)
	}

	return RemoveDuplicatesInPlace(placeholders)
}

package formatter

import (
	"fmt"
	"io"
	"text/tabwriter"

	redwall "github.com/virtualtam/redwall2"
)

func FormatSubmissionAsTab(output io.Writer, submission *redwall.Submission) *tabwriter.Writer {
	writer := tabwriter.NewWriter(output, 0, 4, 2, ' ', 0)

	fmt.Fprintf(writer, "Title\t%s\t\n", submission.Title)
	fmt.Fprintf(writer, "Author\t%s\t\n", submission.User())
	fmt.Fprintf(writer, "Date\t%s\t\n", FormatDateAsUTC(submission.PostedAt))
	fmt.Fprintf(writer, "Post URL\t%s\t\n", submission.PostURL())
	fmt.Fprintf(writer, "Image URL\t%s\t\n", submission.ImageURL)
	fmt.Fprintf(writer, "Image Size\t%d x %d\t\n", submission.ImageWidthPx, submission.ImageHeightPx)
	fmt.Fprintf(writer, "Filename\t%s\t\n", submission.ImageFilename)

	return writer
}

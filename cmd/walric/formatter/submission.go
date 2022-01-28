package formatter

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/virtualtam/walric/submission"
)

// FormatSubmissionAsTab returns a tabwriter.Writer filled with a Submission's
// metadata.
func FormatSubmissionAsTab(output io.Writer, submission *submission.Submission) *tabwriter.Writer {
	writer := tabwriter.NewWriter(output, 0, 4, 2, ' ', 0)

	fmt.Fprintf(writer, "Title\t%s\t\n", submission.Title)
	fmt.Fprintf(writer, "Author\t%s\t\n", submission.User())
	fmt.Fprintf(writer, "Subreddit\t%s\t\n", submission.Subreddit.Name)
	fmt.Fprintf(writer, "Posted At\t%s\t\n", FormatDateAsUTC(submission.PostedAt))
	fmt.Fprintf(writer, "Permalink\t%s\t\n", submission.PermalinkURL())
	fmt.Fprintf(writer, "Image URL\t%s\t\n", submission.ImageURL)
	fmt.Fprintf(writer, "Image Size\t%d x %d\t\n", submission.ImageWidthPx, submission.ImageHeightPx)
	fmt.Fprintf(writer, "Filename\t%s\t\n", submission.ImageFilename)
	fmt.Fprintf(writer, "NSFW\t%t\t\n", submission.ImageNSFW)
	fmt.Fprintf(writer, "Walric ID\t%d\t\n", submission.ID)

	return writer
}

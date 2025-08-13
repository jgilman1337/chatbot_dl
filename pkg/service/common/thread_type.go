package common

// Represents the thread type to export.
type ThreadType int8

const (
	DOCX ThreadType = iota
	Markdown
	PDF
	Log
)

// Emits the extension corresponding to this thread type.
func (t ThreadType) ExtFor() string {
	ext := ""
	switch t {
	case DOCX:
		ext = ".docx"
	case Markdown:
		ext = ".md"
	case PDF:
		ext = ".pdf"
	case Log:
		ext = ".log"
	}
	return ext
}

// Emits the name corresponding to this thread type.
func (t ThreadType) NameFor() string {
	name := ""
	switch t {
	case DOCX:
		name = "DOCX"
	case Markdown:
		name = "Markdown"
	case PDF:
		name = "PDF"
	case Log:
		name = "Logfile"
	}
	return name
}

package common

// Represents an archived thread.
type Thread struct {
	Type     ThreadType //The type of thread this is.
	Filename string     //The filename to use for the thread, excluding the extension.
	Content  []byte     //The content of the thread.
}

// Gets the full filename of the thread, including the extension.
func (t Thread) GetFilename() string {
	return t.Filename + t.Type.ExtFor()
}

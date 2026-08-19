package pdf

type TextExtractor interface {
	Extract(data []byte) (string, error)
}

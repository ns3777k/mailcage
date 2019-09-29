package protocol

type Message struct {
    From string
    To   []string
    Data string
    Helo string
}

package notion

type transcripts []transcript

var DefaultTranscripts = transcripts{
	{Type: "initial", PageTitle: "Improve Coding", PageContent: ""},
	{Type: "human", Prompt: DefaultPrompt},
	{Type: "assistant", Result: "ok"},
}

func (ts transcripts) addHistory(t ...transcript) {
	if len(ts) > transcriptHistorySize {
		ts = append(ts[:2], ts[4:]...)
	}
	for _, v := range t {
		ts = append(ts, v)
	}
}

func (ts transcripts) SetPrompt(p string) {
	DefaultTranscripts[1].Prompt = p
}

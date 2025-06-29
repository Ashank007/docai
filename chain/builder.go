package chain

type ChainBuilder struct {
	embedChain *EmbedChain
	queryChain *QueryChain
}

func NewChainBuilder() *ChainBuilder {
	return &ChainBuilder{}
}

func (b *ChainBuilder) WithEmbedChain(chain *EmbedChain) *ChainBuilder {
	b.embedChain = chain
	return b
}

func (b *ChainBuilder) WithQueryChain(chain *QueryChain) *ChainBuilder {
	b.queryChain = chain
	return b
}

func (b *ChainBuilder) BuildEmbed() Chain {
	return b.embedChain
}

func (b *ChainBuilder) BuildQuery() Chain {
	return b.queryChain
}

package model

type RankOption struct {
	start int64
	stop  int64
}

func (r RankOption) Start() int64 {
	return r.start
}

func (r RankOption) Stop() int64 {
	return r.stop
}

func NewRankOption(start, stop int64) *RankOption {
	return &RankOption{start, stop}
}

func NewAllRankOption() *RankOption {
	return &RankOption{0, -1}
}

func NewMostFrequentOption() *RankOption {
	return &RankOption{-2, -1}
}

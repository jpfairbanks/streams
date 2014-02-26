package sketch

//Countsketch: Implements the Count Sketch Algorithm for frequency estimation in streams
type CountSketch struct {
	hash Hash
	sign NHash
	data Vector
}

//New:TODO determine how much memory to allocate based on n,k,epsilon
func New(n, k int, eps float64) CountSketch {
	return CountSketch{nil, nil, nil}
}

//Insert:Update the sketch by reading a single data element.
func (S CountSketch) Insert(d Datum) error {
	site := S.hash.Apply(d.index)
	S.data[site] += S.sign.Apply(d.index) * d.c
	return nil
}

//Query: Updates q so that the result is the frequency of the index.
func (S CountSketch) Query(q Query) error {
	site := S.hash.Apply(q.index)
	q.result = S.data[site] * S.sign.Apply(q.index)
	return nil
}

//Combine: Updates the current sketch by folding in another sketch T.
//This is necessary for the sketch to be used in parallel reductions.
//If you do not want to overwrite S, then make a copy.
func (S CountSketch) Combine(T Sketch) error {
	target, ok := (T).(CountSketch)
	if !ok {
		return CombinationError("the target T is not the same type as S the host.")
	}
	S.data = S.data.add(target.data)
	return nil
}

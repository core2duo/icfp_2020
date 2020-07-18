package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Atom interface {
	GetName() string
	Evaluate(s Stack) Stack
	Arity() int
	String() string
}

type Stack []Atom

type Partial struct {
	Name string
	Fun  Atom
	Args []Atom
}

func (p *Partial) GetName() string {
	return p.Name
}

func (p *Partial) Arity() int {
	return p.Fun.Arity() - len(p.Args)
}

func (p *Partial) Evaluate(s Stack) Stack {
	if len(s) < 1 {
		log.Panicf("underflow: partial eval: no arg on stack for %s", p.Name)
	}
	arg := s[len(s)-1]
	s = s[0 : len(s)-1]
	p.Args = append(p.Args, arg)
	if p.Arity() == 0 {
		for i := len(p.Args) - 1; i >= 0; i-- {
			s = p.Args[i].Evaluate(s)
		}
		s = p.Fun.Evaluate(s)
	}
	return s
}

func (p *Partial) String() string {
	s := fmt.Sprintf("%s/%d(", p.Name, p.Arity())
	for i := 0; i < len(p.Args); i++ {
		if i > 0 {
			s = s + ", "
		} else {
			s = s + " "
		}
		s = s + p.Args[i].String()
	}
	return s + " )"
}

type Ap struct {
	Fun, Arg Atom
}

func (ap *Ap) GetName() string {
	return ap.Fun.GetName()
}

func (ap *Ap) Arity() int {
	panic("invalid call")
}

func (ap *Ap) Evaluate(s Stack) Stack {
	fun := ap.Fun
	if ap2, ok := fun.(*Ap); ok {
		s = ap2.Evaluate(s)
		fun = s[len(s)-1]
		s = s[0 : len(s)-1]
	}
	arity := fun.Arity()
	if arity == 0 {
		panic("apply arity 0")
	} else if arity == 1 {
		s = ap.Arg.Evaluate(s)
		s = fun.Evaluate(s)
	} else {
		s = ap.Arg.Evaluate(s)
		arg := s[len(s)-1]
		s = s[0 : len(s)-1]
		s = append(s, &Partial{
			Name: fun.GetName(),
			Fun:  fun,
			Args: []Atom{arg},
		})
	}
	return s
}

func (ap *Ap) Construct(s Stack) Stack {
	if len(s) < 2 {
		panic("underflow: ap construct")
	}
	ap.Fun = s[len(s)-1]
	ap.Arg = s[len(s)-2]
	return append(s[0:len(s)-2], ap)
}

func (ap *Ap) String() string {
	return fmt.Sprintf("%s(%s)", ap.Fun.String(), ap.Arg.String())
}

type Pair struct {
	Car, Cdr Atom
}

func (pair *Pair) GetName() string {
	return "pair"
}

var pairStringCounter = 0

func (pair *Pair) String() string {
	pairStringCounter += 1
	var s string
	if pairStringCounter > 10 {
		s = "( ... )"
	} else {
		s = fmt.Sprintf("(%s . %s)", pair.Car.String(), pair.Cdr.String())
	}
	pairStringCounter -= 1
	return s
}

func (pair *Pair) Arity() int {
	return 0
}

func (pair *Pair) Evaluate(s Stack) Stack {
	return append(s, pair)
}

type Cons struct {}

func (cons *Cons) GetName() string {
	return "cons"
}

func (cons *Cons) String() string {
	return "cons"
}

func (cons *Cons) Arity() int {
	return 2
}

func (cons *Cons) Evaluate(s Stack) Stack {
	if len(s) < 2 {
		log.Panicf("underflow: cons eval: %#v", s)
	}
	pair := &Pair{}
	pair.Car = s[len(s)-1]
	pair.Cdr = s[len(s)-2]
	return append(s[0:len(s)-2], pair)
}

type Eq struct {
}

func (eq *Eq) GetName() string {
	return "eq"
}

func (eq *Eq) String() string {
	return "eq"
}

func (eq *Eq) Arity() int {
	return 2
}

func (eq *Eq) Evaluate(s Stack) Stack {
	a := s[len(s)-1].(*Number)
	b := s[len(s)-2].(*Number)
	if a.Value == b.Value {
		s = append(s[0:len(s)-2], &CombK{})
	} else {
		s = append(s[0:len(s)-2], &CombF{})
	}
	return s
}

type Lt struct {
}

func (lt *Lt) GetName() string {
	return "lt"
}

func (lt *Lt) String() string {
	return "lt"
}

func (lt *Lt) Arity() int {
	return 2
}

func (lt *Lt) Evaluate(s Stack) Stack {
	a := s[len(s)-1].(*Number)
	b := s[len(s)-2].(*Number)
	if a.Value < b.Value {
		s = append(s[0:len(s)-2], &CombK{})
	} else {
		s = append(s[0:len(s)-2], &CombF{})
	}
	return s
}

type Add struct {
}

func (add *Add) GetName() string {
	return "add"
}

func (add *Add) String() string {
	return "add"
}

func (add *Add) Arity() int {
	return 2
}

func (add *Add) Evaluate(s Stack) Stack {
	a := s[len(s)-1].(*Number)
	b := s[len(s)-2].(*Number)
	return append(s[0:len(s)-2], &Number{Value: a.Value + b.Value})
}

type Div struct {
}

func (div *Div) GetName() string {
	return "div"
}

func (div *Div) String() string {
	return "div"
}

func (div *Div) Arity() int {
	return 2
}

func (div *Div) Evaluate(s Stack) Stack {
	a := s[len(s)-1].(*Number)
	b := s[len(s)-2].(*Number)
	return append(s[0:len(s)-2], &Number{Value: a.Value / b.Value})
}

type Mul struct {
}

func (mul *Mul) GetName() string {
	return "mul"
}

func (mul *Mul) String() string {
	return "mul"
}

func (mul *Mul) Arity() int {
	return 2
}

func (mul *Mul) Evaluate(s Stack) Stack {
	a := s[len(s)-1].(*Number)
	b := s[len(s)-2].(*Number)
	return append(s[0:len(s)-2], &Number{Value: a.Value * b.Value})
}

type CombK struct {
}

func (combK *CombK) GetName() string {
	return "t"
}

func (combK *CombK) String() string {
	return "t"
}

func (combK *CombK) Arity() int {
	return 2
}

func (combK *CombK) Evaluate(s Stack) Stack {
	x := s[len(s)-1]
	// y := s[len(s)-2]
	return append(s[0:len(s)-2], x)
}

type CombF struct {
}

func (combK *CombF) GetName() string {
	return "f"
}

func (combK *CombF) String() string {
	return "f"
}

func (combK *CombF) Arity() int {
	return 2
}

func (combK *CombF) Evaluate(s Stack) Stack {
	// x := s[len(s)-1]
	y := s[len(s)-2]
	return append(s[0:len(s)-2], y)
}

type CombS struct {
}

func (combS *CombS) GetName() string {
	return "s"
}

func (combS *CombS) String() string {
	return "s"
}

func (combS *CombS) Arity() int {
	return 3
}

func (combS *CombS) Evaluate(s Stack) Stack {
	x := s[len(s)-1]
	y := s[len(s)-2]
	z := s[len(s)-3]
	s = append(s[0:len(s)-3], z)
	s = x.Evaluate(s)
	xz := s[len(s)-1]
	s = append(s, z)
	s = y.Evaluate(s)
	yz := s[len(s)-1]
	s = append(s[0:len(s)-2], yz)
	s = xz.Evaluate(s)
	return s
}

type CombC struct {
}

func (combC *CombC) GetName() string {
	return "c"
}

func (combC *CombC) String() string {
	return "c"
}

func (combC *CombC) Arity() int {
	return 3
}

func (combC *CombC) Evaluate(s Stack) Stack {
	x := s[len(s)-1]
	y := s[len(s)-2]
	z := s[len(s)-3]
	s = append(s[0:len(s)-3], y, z)
	s = x.Evaluate(s)
	xz := s[len(s)-1]
	s = s[0 : len(s)-1]
	s = xz.Evaluate(s)
	return s
}

type CombB struct {
}

func (combB *CombB) GetName() string {
	return "c"
}

func (combB *CombB) String() string {
	return "c"
}

func (combB *CombB) Arity() int {
	return 3
}

func (combB *CombB) Evaluate(s Stack) Stack {
	x := s[len(s)-1]
	y := s[len(s)-2]
	z := s[len(s)-3]
	s = append(s[0:len(s)-3], z)
	s = y.Evaluate(s)
	s = x.Evaluate(s)
	return s
}

type Nil struct{}

func (n *Nil) GetName() string {
	return "nil"
}

func (n *Nil) String() string {
	return "nil"
}

func (n *Nil) Arity() int {
	return 0
}

func (n *Nil) Evaluate(s Stack) Stack {
	s = append(s, n)
	return s
}

type Number struct {
	Value int64
}

func (num *Number) GetName() string {
	return strconv.FormatInt(num.Value, 10)
}

func (num *Number) String() string {
	return strconv.FormatInt(num.Value, 10)
}

func (num *Number) Arity() int {
	return 0
}

func (num *Number) Evaluate(s Stack) Stack {
	return append(s, num)
}

type Neg struct{}

func (neg *Neg) GetName() string {
	return "neg"
}

func (neg *Neg) String() string {
	return "-"
}

func (neg *Neg) Arity() int {
	return 1
}

func (neg *Neg) Evaluate(s Stack) Stack {
	num := s[len(s)-1].(*Number)
	return append(s[0:len(s)-1], &Number{Value: -num.Value})
}

type IsNil struct{}

func (isnil *IsNil) GetName() string {
	return "isnil"
}

func (isnil *IsNil) String() string {
	return "isnil"
}

func (isnil *IsNil) Arity() int {
	return 1
}

func (isnil *IsNil) Evaluate(s Stack) Stack {
	_, ok := s[len(s)-1].(*Nil)
	if ok {
		s = append(s[0:len(s)-1], &CombK{})
	} else {
		s = append(s[0:len(s)-1], &CombF{})
	}
	return s
}

type Car struct{}

func (car *Car) GetName() string {
	return "car"
}

func (car *Car) String() string {
	return "car"
}

func (car *Car) Arity() int {
	return 1
}

func (car *Car) Evaluate(s Stack) Stack {
	pair := s[len(s)-1].(*Pair)
	return append(s[0:len(s)-1], pair.Car)
}

type Cdr struct{}

func (cdr *Cdr) GetName() string {
	return "cdr"
}

func (cdr *Cdr) String() string {
	return "cdr"
}

func (cdr *Cdr) Arity() int {
	return 1
}

func (cdr *Cdr) Evaluate(s Stack) Stack {
	pair := s[len(s)-1].(*Pair)
	return append(s[0:len(s)-1], pair.Cdr)
}

type CombI struct{}

func (neg *CombI) GetName() string {
	return "i"
}

func (neg *CombI) String() string {
	return "i"
}

func (neg *CombI) Arity() int {
	return 1
}

func (neg *CombI) Evaluate(s Stack) Stack {
	return s
}

type Ref struct {
	Name  string
	Words []string
}

func (ref *Ref) GetName() string {
	return ref.Name
}

func (ref *Ref) String() string {
	return ref.Name
}

func (ref *Ref) Arity() int {
	panic("not implemented")
}

func (ref *Ref) Evaluate(s Stack) Stack {
	return env[ref.Name][0].Evaluate(s)
}

var env map[string][]Atom

func parse(lets map[string][]string, v []string) []Atom {
	s := make(Stack, 0, len(v))
	for i := len(v) - 1; i >= 0; i-- {
		word := v[i]
		switch word {
		case "ap":
			ap := &Ap{}
			s = ap.Construct(s)
		case "nil":
			s = append(s, &Nil{})
		case "cons":
			s = append(s, &Cons{})
		case "neg":
			s = append(s, &Neg{})
		case "i":
			s = append(s, &CombI{})
		case "c":
			s = append(s, &CombC{})
		case "b":
			s = append(s, &CombB{})
		case "t":
			s = append(s, &CombK{})
		case "f":
			s = append(s, &CombF{})
		case "car":
			s = append(s, &Car{})
		case "cdr":
			s = append(s, &Cdr{})
		case "isnil":
			s = append(s, &IsNil{})
		case "eq":
			s = append(s, &Eq{})
		case "lt":
			s = append(s, &Lt{})
		case "s":
			s = append(s, &CombS{})
		case "add":
			s = append(s, &Add{})
		case "mul":
			s = append(s, &Mul{})
		case "div":
			s = append(s, &Div{})
		default:
			_, ok := lets[word]
			if ok {
				s = append(s, &Ref{Name: word, Words: lets[word]})
			} else {
				num, err := strconv.ParseInt(word, 10, 64)
				if err != nil {
					log.Panicf("cannot parse: %s", word)
				}
				s = append(s, &Number{num})
			}
		}
	}
	return s
}

func show(a Atom) string {
	if ap, ok := a.(*Ap); ok {
		args := []string{}
		ap2 := ap
		fun := ap.Fun
		for {
			ap2, ok = fun.(*Ap)
			if ok {
				args = append([]string{show(ap2.Arg) + ", "}, args...)
				fun = ap2.Fun
			} else {
				break
			}
		}
		return fmt.Sprintf("%s (%s%s)", show(fun), strings.Join(args, ""), show(ap.Arg))
	} else {
		return a.GetName()
	}
}

func load() {
	f, err := os.Open("galaxy.txt")
	if err != nil {
		log.Panicf("cannot open galaxy.txt: %s", err)
	}

	env = make(map[string][]Atom)
	lets := make(map[string][]string)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		kv := strings.SplitN(line, " = ", 2)
		if len(kv) < 2 {
			log.Printf("bad line: %q", line)
			continue
		}
		lets[kv[0]] = strings.Fields(kv[1])
	}
	if err := scanner.Err(); err != nil {
		log.Panicf("scanner has failed: %s", err)
	}
	f.Close()

	for k, v := range lets {
		env[k] = parse(lets, v)
		if len(env[k]) != 1 {
			panic("bad parse result length")
		}
	}

	s := make(Stack, 0)
	log.Println(show(env["galaxy"][0]))
	s = env["galaxy"][0].Evaluate(s)
	log.Println(s)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	load()
}

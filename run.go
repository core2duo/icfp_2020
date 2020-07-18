package main

import (
	"bufio"
	"flag"
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
		if ap2, ok := ap.Arg.(*Ap); ok {
			s = ap2.Evaluate(s)
		} else if ap.Arg.Arity() == 0 {
			s = ap.Arg.Evaluate(s)
		} else {
			s = append(s, &Partial{
				Name: ap.Arg.GetName(),
				Fun:  ap.Arg,
			})
		}
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

type Cons struct{}

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

type Mod struct{}

func (mod *Mod) GetName() string {
	return "mod"
}

func (mod *Mod) String() string {
	return "mod"
}

func (mod *Mod) Arity() int {
	return 1
}

func (mod *Mod) Evaluate(s Stack) Stack {
	return append(s[0:len(s)-1], &Linear{Data: modulate(s[len(s)-1])})
}

type Linear struct {
	Data string
}

func (linear *Linear) GetName() string {
	return "modulated-data"
}

func (linear *Linear) String() string {
	return "[" + linear.Data + "]"
}

func (linear *Linear) Arity() int {
	return 0
}

func (linear *Linear) Evaluate(s Stack) Stack {
	return append(s, linear)
}

type Dem struct{}

func (dem *Dem) GetName() string {
	return "dem"
}

func (dem *Dem) String() string {
	return "dem"
}

func (dem *Dem) Arity() int {
	return 1
}

func (dem *Dem) Evaluate(s Stack) Stack {
	arg := s[len(s)-1].(*Linear)
	atom, rest := demodulate(arg.Data)
	if rest != "" {
		log.Fatalf("cannot demodulate %q: unexpected trailer: %q", arg.Data, rest)
	}
	return append(s[0:len(s)-1], atom)
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
	return 1 // FIXME: shouldn't we actually calculate the arity of refs?
}

func (ref *Ref) Evaluate(s Stack) Stack {
	return env[ref.Name].Evaluate(s)
}

var env map[string]Atom

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
		case "mod":
			s = append(s, &Mod{})
		case "dem":
			s = append(s, &Dem{})
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
	f, err := os.Open(inputFile)
	if err != nil {
		log.Panicf("cannot open %q: %s", inputFile, err)
	}

	env = make(map[string]Atom)
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

	lets["main"] = strings.Fields(initExpr)

	for k, v := range lets {
		result := parse(lets, v)
		if len(result) != 1 {
			panic("bad parse result length")
		}
		env[k] = result[0]
	}

	s := make(Stack, 0)
	log.Printf("   %s", show(env["main"]))
	s = env["main"].Evaluate(s)
	log.Printf("=> %s", s[0])
}

func modulate(a Atom) string {
	switch v := a.(type) {
	case *Number:
		prefix := "01"
		val := v.Value
		if v.Value < 0 {
			prefix = "10"
			val = -v.Value
		}
		if val == 0 {
			return prefix + "0"
		}
		bstr := fmt.Sprintf("%b", val)
		for len(bstr)%4 > 0 {
			bstr = "0" + bstr
		}
		for i := 0; i < len(bstr); i += 4 {
			prefix = prefix + "1"
		}
		return prefix + "0" + bstr

	case *Nil:
		return "00"

	case *Pair:
		return "11" + modulate(v.Car) + modulate(v.Cdr)
	}

	log.Fatalf("cannot modulate %T: %#v", a, a)
	return ""
}

func demodulate(s string) (Atom, string) {
	sign := int64(1)

	switch s[0:2] {
	case "00":
		return &Nil{}, s[2:]

	case "10":
		sign = -1
		fallthrough

	case "01":
		bits, i := 0, 0
		for i = 2; i < len(s) && s[i] == '1'; i++ {
			bits += 4
		}
		if i+bits+1 > len(s) || s[i] != '0' {
			log.Fatalf("unexpected end of string: %q", s)
		}
		i += 1
		v, err := strconv.ParseInt(s[i:i+bits], 2, 64)
		if err != nil {
			log.Fatalf("cannot parse %q: %s", s[i:i+bits], err)
		}
		return &Number{Value: sign * v}, s[i+bits:]

	case "11":
		car, rest := demodulate(s[2:])
		cdr, rest := demodulate(rest)
		return &Pair{Car: car, Cdr: cdr}, rest

	default:
		log.Fatalf("bad tag: %q", s[0:2])
	}

	return nil, ""
}

var initExpr string
var inputFile string

func init() {
	flag.StringVar(&initExpr, "expr", "galaxy", "Expression to evaluate")
	flag.StringVar(&inputFile, "in", "galaxy.txt", "Input file")
	flag.Parse()
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	load()
}

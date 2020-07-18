class Atom(object):
    def __init__(self, name=''):
        self.name = name

    def __repr__(self):
        return self.name


class Number(Atom):
    def __init__(self, value):
        self.name = str(value)
        self.value = int(value)

    def eval(self):
        return self.value


class TrueAtom(Atom):
    arity = 2

    def __init__(self, name=''):
        self.name = 't'

    def load(self, bt, tf):
        self.bt = bt
        self.bf = bf

    def eval(self):
        return self.bt

    def __repr__(self):
        return '{!r} {!r} {!r}'.format(self.name, self.bt, self.bf)


class FalseAtom(TrueAtom):
    def eval(self):
        return self.bf

    def __init__(self, name=''):
        self.name = 'f'


class Add(Atom):
    arity = 2

    def __init__(self, name=''):
        self.name = 'add'

    def load(self, a, b):
        self.a = a
        self.b = b

    def eval(self):
        return self.a.eval() + self.b.eval()

    def __repr__(self):
        return '{!r} + {!r}'.format(self.a, self.b)


class Mul(Atom):
    arity = 2

    def __init__(self, name=''):
        self.name = 'mul'

    def load(self, a, b):
        self.a = a
        self.b = b

    def eval(self):
        return self.a.eval() * self.b.eval()

    def __repr__(self):
        return '{!r} * {!r}'.format(self.a, self.b)


class Div(Atom):
    arity = 2

    def __init__(self, name=''):
        self.name = 'div'

    def load(self, a, b):
        self.a = a
        self.b = b

    def eval(self):
        return self.a.eval() / self.b.eval()

    def __repr__(self):
        return '{!r} / {!r}'.format(self.a, self.b)


class Eq(Atom):
    arity = 2

    def __init__(self, name=''):
        self.name = 'eq'

    def load(self, a, b):
        self.a = a
        self.b = b

    def eval(self):
        if self.a.eval() == self.b.eval():
            return TrueAtom('t')
        else:
            return FalseAtom('f')

    def __repr__(self):
        return '{!r} == {!r}'.format(self.a, self.b)


class Lt(Atom):
    arity = 2

    def __init__(self, name=''):
        self.name = 'lt'

    def load(self, a, b):
        self.a = a
        self.b = b

    def eval(self):
        if self.a.eval() < self.b.eval():
            return TrueAtom('t')
        else:
            return FalseAtom('f')

    def __repr__(self):
        return '{!r} < {!r}'.format(self.a, self.b)


class I(Atom):
    arity = 1

    def __init__(self, name=''):
        self.name = 'i'

    def load(self, value):
        self.value = value

    def eval(self):
        return self.value.eval()

    def __repr__(self):
        return '{!r}'.format(self.value)


class Nil(Atom):
    arity = 0

    def __init__(self, name=''):
        self.name = 'nil'

    def load(self):
        pass

    def eval(self):
        return None


class Neg(Atom):
    arity = 1

    def __init__(self, name=''):
        self.name = 'neg'

    def load(self, value):
        self.value = value

    def eval(self):
        return -self.value.eval()

    def __repr__(self):
        return '- {!r}'.format(self.value)


class Thunk(object):
    def __init__(self, kind, args=None):
        self.kind = kind
        self.args = args or []

    def __repr__(self):
        return '{}/{} ({})'.format(self.kind, len(self.args), ', '.join(map(repr, self.args)))

    def push(self, value):
        self.args.append(value)
        if self.kind.arity == len(self.args):
            value = self.kind()
            value.load(*self.args)
            return value
        else:
            return self


class Parser(object):
    def __init__(self, env):
        self.env = env

    def parse(self, value):
        self.stack = []
        words = reversed(value.strip().split())
        for word in words:
            if word == 'ap':
                fun = self.stack.pop()
                assert len(self.stack) > 0, 'stack underflow for {!r}'.format(fun)
                arg = self.stack.pop()
                if isinstance(fun, type) and issubclass(fun, Atom):
                    self.stack.append(Thunk(fun).push(arg))
                elif isinstance(fun, Thunk):
                    self.stack.append(fun.push(arg))
                else:
                    assert RuntimeError('unexpected fun type for ap: {!r}'.format(fun))
            elif word in SPECIALS:
                special = SPECIALS[word]
                if special.arity == 0:
                    self.stack.append(special(word))
                else:
                    self.stack.append(special)
            elif word in env:
                self.stack.append(env[word])
            elif word.startswith('-') and word[1:].isnumeric():
                self.stack.append(Number(word))
            elif word.isnumeric():
                self.stack.append(Number(word))
            else:
                raise NameError("name '{}' is not defined".format(word))
        assert len(self.stack) == 1
        return self.stack[0]


class Cons(Atom):
    arity = 2
    _rec = 0

    def __init__(self, name=''):
        self.name = 'cons'

    def load(self, car, cdr):
        self.car = car
        self.cdr = cdr

    def eval(self):
        raise TypeError("'Cons' object is not callable")

    def __repr__(self):
        if Cons._rec > 5:
            return '...'
        Cons._rec = Cons._rec + 1
        if isinstance(self.cdr, Cons):
            value = '[{!r}'.format(self.car)
            cdr = self.cdr
            while isinstance(cdr, Cons):
                value = value + ' {!r}'.format(cdr.car)
                cdr = cdr.cdr
            if isinstance(cdr, Nil):
                value = value + ']'
            else:
                value = value + ' . {!r}'.format(cdr)
        else:
            value = '[{!r} . {!r}]'.format(self.car, self.cdr)
        Cons._rec = Cons._rec - 1
        return value


class C(Atom):
    arity = 3

    def __init__(self, name=''):
        self.name = 'C'

    def load(self, x, y, z):
        self.x = x
        self.y = y
        self.z = z

    def eval(self):
        x = self.x.eval()
        z = self.z.eval()
        y = self.y.eval()

    def __repr__(self):
        return '{!r} ({!r}) ({!r})'.format(self.x, self.z, self.y)


class B(Atom):
    arity = 3

    def __init__(self, name=''):
        self.name = 'B'

    def load(self, x, y, z):
        self.x = x
        self.y = y
        self.z = z

    def eval(self):
        x = self.x.eval()
        y = self.y.eval()
        z = self.z.eval()

    def __repr__(self):
        return '{!r} ({!r} ({!r}))'.format(self.x, self.y, self.z)


class S(Atom):
    arity = 3

    def __init__(self, name=''):
        self.name = 'S'

    def load(self, x, y, z):
        self.x = x
        self.y = y
        self.z = z

    def eval(self):
        x = self.x.eval()
        y = self.y.eval()
        z = self.z.eval()

    def __repr__(self):
        return 'S({!r}, {!r}, {!r})'.format(self.x, self.y, self.z)


class Car(Atom):
    arity = 1

    def __init__(self, name=''):
        self.name = 'car'

    def load(self, cons):
        self.cons = cons

    def eval(self):
        cons = self.cons.eval()
        assert isinstance(cons, Cons)
        return cons.car

    def __repr__(self):
        return 'car ({!r})'.format(self.cons)


class Cdr(Atom):
    arity = 1

    def __init__(self, name=''):
        self.name = 'cdr'

    def load(self, cons):
        self.cons = cons

    def eval(self):
        cons = self.cons.eval()
        assert isinstance(cons, Cons)
        return cons.cdr

    def __repr__(self):
        return 'cdr ({!r})'.format(self.cons)


class IsNil(Atom):
    arity = 1

    def __init__(self, name=''):
        self.name = 'isnil'

    def load(self, cons):
        self.cons = cons

    def eval(self):
        cons = self.cons.eval()
        if isinstance(cons, Cons):
            return SPECIALS['t']('t')
        else:
            return SPECIALS['f']('f')

    def __repr__(self):
        return 'isnil ({!r})'.format(self.cons)


SPECIALS = {
             'cons': Cons,
             'nil': Nil,
             'neg': Neg,
             'c': C,
             'b': B,
             'car': Car,
             'cdr': Cdr,
             't': TrueAtom,
             'f': FalseAtom,
             'isnil': IsNil,
             's': S,
             'add': Add,
             'mul': Mul,
             'eq': Eq,
             'div': Div,
             'lt': Lt,
             'i': I,
           }

env = {
      }

with open('galaxy.txt') as f:
    for line in f:
        ident, value = line.strip().split(' = ', 2)
        atom = Atom(ident)
        atom.value = value
        env[ident] = atom


for ident, atom in env.items():
    if type(atom) is Atom:
        parser = Parser(env)
        try:
            code = parser.parse(atom.value)
        except:
            print('Error parsing', ident)
            raise
        env[ident] = code
        print('{} = {!r}'.format(ident, code))

class Atom(object):
    def __init__(self, name=''):
        self.name = name

    def __repr__(self):
        return self.name


class Number(Atom):
    def __init__(self, value):
        self.name = 'Number<{}>'.format(value)
        self.value = int(value)

    def eval(self):
        return self.value, words


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
        self.value = None

    def load(self, value):
        self.value = value

    def eval(self):
        return None

    def __repr__(self):
        return '- {!r}'.format(self.value)


class Thunk(object):
    def __init__(self, kind, args=None):
        self.kind = kind
        self.args = args or []

    def __repr__(self):
        return '{}/{}({})'.format(self.kind, len(self.args), ', '.join(map(repr, self.args)))

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
                assert len(self.stack) >= 2, 'stack underflow: {!r}'.format(self.stack)
                fun = self.stack.pop()
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
        value = '<{!r} . {!r}>'.format(self.car, self.cdr)
        Cons._rec = Cons._rec - 1
        return value


SPECIALS = {
             'cons': Cons,
             'nil': Nil,
             'neg': Neg,
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
        code = parser.parse(atom.value)
        env[ident] = code
        print('{} = {!r}'.format(ident, code))

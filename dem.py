def assert_equal(a, b):
    assert a == b, '{!r} == {!r}'.format(a, b)


def dem(str):
    """
    Demodulates a string of 0s and 1s.
    Returns the demodulated value and the rest of the string.
    """

    if str.startswith('10'):  # negative number
        num, rest = dem('01' + str[2:])
        return (-num, rest)
    elif str.startswith('01'):  # positive number
        i = 2
        bits = 0
        num = 0
        while i < len(str) and str[i] == '1':
            i = i + 1
            bits = bits + 1
        assert str[i] == '0'
        i = i + 1
        rembits = bits
        while i < len(str) and rembits > 0:
            for _ in range(4):
                num = num * 2 + int(str[i])
                i = i + 1
            rembits -= 1
        assert_equal(rembits, 0)
        return num, str[i:]
    elif str.startswith('00'):
        return [], str[2:]
    elif str.startswith('11'):
        car, rest = dem(str[2:])
        cdr, rest = dem(rest)
        return [car, cdr], rest

assert_equal(dem('010'), (0, ''))
assert_equal(dem('100'), (-0, ''))
assert_equal(dem('00'), ([], ''))
assert_equal(dem('110000'), ([[], []], ''))
assert_equal(dem('1101000'), ([0, []], ''))
assert_equal(dem('110110000101100010'), ([1, 2], ''))

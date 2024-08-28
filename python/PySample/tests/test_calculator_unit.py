import unittest
from PySample.pysample import calculator


class TestStringMethods(unittest.TestCase):

    def test_addition2(self):
        self.assertEqual(4, calculator.add(2, 2))

    def test_addition3(self):
        self.assertEqual(2, calculator.subtract(4, 2))


if __name__ == '__main__':
    unittest.main()

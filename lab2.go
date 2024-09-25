/*
1. Написать программу, которая определяет, является ли введенное пользователем число четным или нечетным.
2. Реализовать функцию, которая принимает число и возвращает "Positive", "Negative" или "Zero".
3. Написать программу, которая выводит все числа от 1 до 10 с помощью цикла for.
4. Написать функцию, которая принимает строку и возвращает ее длину.
5. Создать структуру Rectangle и реализовать метод для вычисления площади прямоугольника.
6. Написать функцию, которая принимает два целых числа и возвращает их среднее значение.
*/

package main

import (
	"fmt"
)

func main() {
	var (
		n    int
		num1 int
		num2 int
		s    string
		r    Rectangle
		a    int
		b    int
	)
	for {
		fmt.Println("Введите номер задания (0 - выход)")
		fmt.Scanln(&n)
		if n == 0 {
			break
		}
		switch n {
		case 1:
			fmt.Println("Введите число")
			fmt.Scanln(&num1)
			if num1%2 == 0 {
				fmt.Println("Число четное")
			} else {
				fmt.Println("Число нечетное")
			}
		case 2:
			fmt.Println("Введите число")
			fmt.Scanln(&num1)
			PNZ(num1)
		case 3:
			for i := 1; i <= 10; i++ {
				fmt.Println(i)
			}
		case 4:
			fmt.Println("Введите строку (без пробелов)")
			fmt.Scanln(&s)
			Length(s)
		case 5:
			fmt.Println("Введите длину и ширину прямоугольника")
			fmt.Scanln(&a)
			fmt.Scanln(&b)
			r = Rectangle{a, b}
			Square(r)
		case 6:
			fmt.Println("Введите два числа")
			fmt.Scanln(&num1)
			fmt.Scanln(&num2)
			Avg(num1, num2)
		default:
			fmt.Println("Ошибка! Нет задания с таким номером")
		}
	}

}

func PNZ(num int) {
	if num == 0 {
		fmt.Println("Zero")
	} else if num < 0 {
		fmt.Println("Negative")
	} else {
		fmt.Println("Positive")
	}
}

func Length(s string) {
	var l = len([]rune(s))
	fmt.Println("Количество символов в строке: ", l)
}

type Rectangle struct {
	a int
	b int
}

func Square(r Rectangle) {
	var s = r.a * r.b
	fmt.Println("Площадь прямоугольника ", s)
}

func Avg(num1, num2 int) {
	var avg int = (num1 + num2) / 2
	fmt.Println("Среднее арифметическое = ", avg)
}

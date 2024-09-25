/*
1. Написать программу, которая выводит текущее время и дату.
2. Создать переменные различных типов (int, float64, string, bool) и вывести их на экран.
3. Использовать краткую форму объявления переменных для создания и вывода переменных.
4. Написать программу для выполнения арифметических операций с двумя целыми числами и выводом результатов.
5. Реализовать функцию для вычисления суммы и разности двух чисел с плавающей запятой.
6. Написать программу, которая вычисляет среднее значение трех чисел.
*/

package main

import (
	"fmt"
	"time"
)

func lab1() {
	fmt.Println("Введите номер задания (0 - выход)")
	var (
		n    int
		num1 int
		num2 int
		num3 float64
		num4 float64
		num5 int
	)
	fmt.Scanln(&n)
	for {
		if n == 0 {
			break
		}
		switch n {
		case 1:
			fmt.Println("current date + time:")
			fmt.Println(time.Now().Local().Format(time.RFC1123))
		case 2:
			var (
				int_num   int     = 12
				float_num float64 = 11.38
				str       string  = "golang"
				b         bool    = true
			)
			fmt.Println(int_num, float_num, str, b)
		case 3:
			Num_Int, Num_Float64, Str, B := 42, 12.8, "go", false
			fmt.Println(Num_Int, Num_Float64, Str, B)
		case 4:
			fmt.Println("Введите два целых числа")
			fmt.Scanln(&num1)
			fmt.Scanln(&num2)
			IntOperations(num1, num2)
		case 5:
			fmt.Println("Введите два числа с плавающей запятой")
			fmt.Scanln(&num3)
			fmt.Scanln(&num4)
			fmt.Println("")
			FlOperations(num3, num4)
		case 6:
			fmt.Println("Введите три целых числа")
			fmt.Scanln(&num1)
			fmt.Scanln(&num2)
			fmt.Scanln(&num5)
			Average(num1, num2, num5)
		default:
			fmt.Println("Ошибка! Нет задания с таким номером")
		}
	}

}

func IntOperations(num1, num2 int) {
	var (
		sum  int = num1 + num2
		ras  int = num1 - num2
		mult int = num1 * num2
	)
	fmt.Println("Сумма чисел = ", sum)
	fmt.Println("Разность чисел = ", ras)
	fmt.Println("Произведение чисел = ", mult)
	if num2 == 0 {
		fmt.Println("Ошибка! Нельзя делить на ноль!")
	} else {
		var del int = num1 / num2
		fmt.Println("Частное чисел = ", del)
	}
}

func FlOperations(num1, num2 float64) {
	var (
		sum float64 = num1 + num2
		ras float64 = num1 - num2
	)
	fmt.Println("Сумма чисел = ", sum)
	fmt.Println("Разность чисел = ", ras)
}

func Average(num1, num2, num3 int) {
	var avg int = (num1 + num2 + num3) / 3
	fmt.Println("Среднее арифметическое = ", avg)
}

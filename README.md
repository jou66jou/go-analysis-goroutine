原始問題出處：  
<https://stackoverflow.com/questions/35153010/goroutines-always-execute-last-in-first-out>  
<https://studygolang.com/topics/9608>  
***

問題範例code:  
```go

func c(i int){
  fmt.Println(i)
}


func main() {
   runtime.GOMAXPROCS(1)

   for i:=0; i< 11; i++ {
       go c(i) // print依序為10 0 1 2 3 4 5 6 7 8 9
   }

  fmt.Scanln()
}

```  
## 為什麼是10先被打印出來？  
首先要先知道goroutine隊列執行本身就可能是無序的，其原理可去理解Golang的GMP模型  
單核心在數字大時就不會總是最後一個被優先print出來、也不會總是照順序  
接著我們用trace tool做三個實驗，看看單核心在什麼情況下`最後一個數字不會被先print`  
***
### * 實驗一：for 0~999 單核心

結果：執行後第一個被print的為999  
![螢幕快照 2019-07-18 下午5.02.58.png](https://static.studygolang.com/190718/0eb639533c172c431ff4654c3426bf85.png)

1.先看trace總覽，goroutines在runtime.main執行結束後開始下降（青藍色最高峰），也就是1000個G開始被調用消耗掉  
  
![螢幕快照 2019-07-18 下午4.49.47.png](https://static.studygolang.com/190718/df5e8dc9e441fb66f58403b24dcd18bd.png)  
  
2.接著看看第一個被調用的G，就是print數字999的goroutine，可以看到Start Stack Trace有成功執行syscall.write，將999打印到螢幕上
![螢幕快照 2019-07-18 下午4.39.00.png](https://static.studygolang.com/190718/1b7d22950024bf86f24b09e70f42759c.png)

***
### * 實驗二：for 0~9999 單核心
結果：執行後第一個不為9999  
![螢幕快照 2019-07-18 下午5.14.26.png](https://static.studygolang.com/190718/7dfb0c61cb35c1baa72800b11c3c80ef.png)  
  
1.與實驗一相同，goroutines在runtime.main執行結束後才開始下降，但可以注意到GC出現了！
![螢幕快照 2019-07-18 下午5.11.09.png](https://static.studygolang.com/190718/0da9ff237af4387355508a1e5f4b1e3b.png)
  
2.來到第一個被調用的Ｇ，可以發現應該就是要執行print數字9999的goroutine，`卻沒有執行syscall.write`，而是進入GC的程序  
![螢幕快照 2019-07-18 下午4.38.34.png](https://static.studygolang.com/190718/5d656f61ea7d5d195dc99ec15c882bbb.png)

***
### * 實驗三：for 0~999 多核心
結果：執行後亂序出現  
![螢幕快照 2019-07-18 下午5.20.41.png](https://static.studygolang.com/190718/54c7891880f7ea0fb632eb4cb1f8a4c5.png)
  
1.從trace可以發現goroutines***不會再等待*** runtime.main執行完畢才開始被調用  
![螢幕快照 2019-07-18 下午5.22.00.png](https://static.studygolang.com/190718/ab8646a116017a8322b8137ec6085468.png)  
  
***
### * 實驗總結：
我們可以發現在單核心情況下，runtime.main的過程中G不斷在被創建，然而要等到runtime.main結束才開始消耗G。  
* 問題一： 為什麼第一個被調用的G是for迴圈中最後一個999呢？  
`解釋`  
stackoverflow也有類似問題被問過(註一)，但回答依舊是goroutine本身隊列無序，因此我們只能大略推論其邏輯是為了***效率***，也就是最後一個goroutine被創建後就直接消耗掉，因此才會出現第一個被執行的G—是for迴圈中最後一個被創建的G。  
* 問題二：同樣是單核心情況，為什麼迴圈數字變大，第一個print的不是最後一個數字呢？  
`解釋`  
從實驗二可以發現是***GC從中作梗***，第一個被調用的G仍然是最後被創建的G，但是因為GC的作用下syscall.write被延後執行。  
* 問題三：單核心跟多核心有什麼差異？  
`解釋`  
實驗三可以簡潔明瞭發現問題在runtime.main創建G時，`其他CPU可以同時執行G`，就不會有最後創建的G第一個被消耗情況。  


註一：<https://stackoverflow.com/questions/35153010/goroutines-always-execute-last-in-first-out>

後記：  
解答兄弟問題，為我快樂之本  
此篇同時更新到個人的github，若喜歡請幫忙給個星星start，有任何觀念錯誤也請多多指教～  
<https://github.com/jou66jou/go-analysis-goroutine>
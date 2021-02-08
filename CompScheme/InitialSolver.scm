#lang racket

(current-directory)


(define (readTableau fileIn)  
  (let ((sL (map (lambda s (string-split (car s))) (file->lines fileIn))))
    (map (lambda (L)
           (map (lambda (s)
                  (if (eqv? (string->number s) #f)
                      s
                      (string->number s))) L)) sL)))
    

(define (writeTableau tb fileOut)
  (if (eqv? tb '())
      #t
      (begin (display-lines-to-file (car tb) fileOut #:separator #\space #:exists 'append)
             (display-to-file #\newline fileOut #:exists 'append)
             (writeTableau (cdr tb) fileOut))))
                             


;Gets the number of rows in a grid
(define (getNumberOfRows l) 
  (if (null? l) 0 
    (+ 1 (getNumberOfRows (cdr l)))
  )
)

;Gets the number of col in a grid
(define (getNumberOfCols l)
  (let ((firstrow (car l)))
    (getNumberOfRows firstrow)
  )
)

(define (getCostGrid l)
  (let* ((newGrid (cdr l)) (firstrow (car newGrid)))
    (if (null? (cdr newGrid)) '()
      (cons (cdr firstrow) (getCostGrid newGrid))
    )
  )
)

;(define costGrid (getCostGrid grid))

(define (getDemandList l)
  (if (null? (cdr l)) (cdr(car l))
    (getDemandList (cdr l))
  )
)



;Apply function f to element of 1-D list in position col and return the modified list
;Assumes column is within range
(define (applyFunctionToElement l col f)
  (if (= col 0) (cons (f (car l)) (cdr l))
    (cons (car l) (applyFunctionToElement (cdr l) (- col 1) f))
  )
)

;Apply function f to element of 2-d list in position col,row and returns the modified list
;Assumes row is within range
(define (applyFunctionToCell g row col f)
  (if (= row 0) (cons (applyFunctionToElement (car g) col f) (cdr g))
    (cons (car g) (applyFunctionToCell (cdr g) (- row 1) col f))
  )
)

;Gets element at position col in the 1-D list l
(define (getValueOfCol l col)
  (if (= col 0) (car l)
    (getValueOfCol (cdr l) (- col 1))
  )
)

;Gets element at position row col in the 2-D list g
(define (getValueOfCell g row col)
  (if (= row 0) (getValueOfCol (car g) col)
    (getValueOfCell (cdr g) (- row 1) col)
  )
)

;Gets column col of a grid as a list
(define (getCol g col)
  (if (null? g) '()
  (cons (getValueOfCol (car g) col) (getCol (cdr g) col))))

(define (getSupplyList grid costGrid)
  (let ((num (getNumberOfCols grid)))
    (getCol costGrid (- num 2))))


;Gets the index of the min element in column col in grid g
(define (getIndexOfMin g col) (let
  ((lcol (getCol g col)))
  (do ((i 0 (+ i 1)) (minIndex 0)(min (car lcol)) (remaining lcol (cdr remaining)))
    ((null? remaining) minIndex)
    (unless (>= (car remaining) min)(set! min (car remaining))(set! minIndex i) ))
  )
)

;Returns grid of 0s with size rowsxcol
(define (getEmptyList size)
  (if (= size 1) (list 0)
    (cons 0 (getEmptyList (- size 1)))))

;Returns grid of size "rows"X"col" with rows "row" ("" to distinguish between parameters and words to be taken literally)
(define (getEmptyGrid rows col row)
  (if (= 1 rows) (list row)
    (cons row (getEmptyGrid (- rows 1) col row))
  ))

;Auxillary function that creates an empty row of size col and then calls the getEmptyGrid function with that row (quick fix for memory efficieny)
(define (getEmptyGridAux rows col)
  (getEmptyGrid rows col (getEmptyList col)))


(define (initialSolution costGrid demandList supplyList numRows numColsCost) 
  (do ((i 0) (grid costGrid) (demands demandList) (supplies supplyList)(finalGrid (getEmptyGridAux numRows numColsCost)))
    ((= i numColsCost) finalGrid)
     (let* ((demand (getValueOfCol demands i))(minAvailRow (getIndexOfMin grid i))(supply (getValueOfCol supplies minAvailRow)) (toAdd (min supply demand)) (increment (lambda (x) (+ x toAdd)))
    (decrement (lambda(x) (- x toAdd))) (explode (lambda(x) 100000000)))
    (cond ((not(= demand 0))
    (set! supplies (applyFunctionToElement supplies minAvailRow decrement))
    (set! demands (applyFunctionToElement demands i decrement))
    (set! grid (applyFunctionToCell grid minAvailRow i explode))
    (set! finalGrid (applyFunctionToCell finalGrid minAvailRow i increment)))
    (else (set! i (+ i 1))) 
    )
  )
))

;Next couple of definitions are to modify the grid that will be written to the text file to the required format

;Appends column passed as a list col to end of grid g 
(define (appendColumn g col) 
  (if (null? g) '()
    (cons (append (car g) (list (car col))) (appendColumn (cdr g) (cdr col)))))

(define (addColumnFirst g col)
  (if (null? g) '()
    (cons (cons (car col) (car g)) (addColumnFirst (cdr g) (cdr col)))))

(define (removeLastElement l)
  (if (null? (cdr l)) '()
    (cons (car l) (removeLastElement (cdr l)))))

(define (gridHeader grid) (car grid))

(define (gridSideBar grid)
  (removeLastElement (cdr (getCol grid 0))))

(define (gridTail demandList) 
  (cons "DEMAND" demandList))

(define (formattedSolution solvedGrid gridHeader gridSideBar supplyList gridTail)
  (append (cons gridHeader (appendColumn (addColumnFirst solvedGrid gridSideBar) supplyList)) (list gridTail)))


(define (minimumCell fileIn fileOut)
  (let* ((g (readTableau fileIn))(grid g) (costGrid (getCostGrid grid)) (demandList (getDemandList grid)) (supplyList (getSupplyList grid costGrid)) (numColsCost (- (getNumberOfCols costGrid) 1)) (numRows (getNumberOfRows costGrid))) (writeTableau (formattedSolution (initialSolution costGrid demandList supplyList numRows numColsCost) (gridHeader grid) (gridSideBar grid) supplyList (gridTail demandList)
) fileOut)))




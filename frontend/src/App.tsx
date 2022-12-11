import { useCallback, useEffect, useRef, useState } from "react";
import "./App.css";
import { GenerateBoard, GetBoard, Solve } from "../wailsjs/go/main/App";
import {
  Text,
  Button,
  ChakraProvider,
  Grid,
  GridItem,
  Heading,
} from "@chakra-ui/react";
import { Cycle, CycleState, motion } from "framer-motion";

const EMPTY = 0;

enum Status {
  FAILURE = 0,
  SUCCESS = 1,
  CUTOFF = 2,
  TIME_EXCEEDED = 3,
}

type SolveData = {
  Status: Status;
  Iterations: number[][];
  IterationCount: number;
  //time elapsed in milliseconds
  TimeElapsed: number;
};

type Cell = {
  value: number;
};

const isSolved = (board: number[]) : boolean => {
  if(!board) return false
  for (let index = 0; index < board.length; index++) {
    if(index === board.length -1 && board[index] === 0) return true;
    if (board[index] !== index+1) {
      return false
    }
  }
  return false
}

function App() {
  // const [emptyIndex, setEmptyIndex] = useState<number | undefined>(undefined);
  const [isSolving, setSolving] = useState(false);
  const [isAnimating, setAnimating] = useState(false);
  const [solveData, setSolveData] = useState<SolveData | undefined>();
  const [boards, setBoards] = useState<number[][] | undefined>();

  // useEffect(() => {
  //   if (boards && boards.length > 0) {
  //     setEmptyIndex(boards[0].findIndex((i) => i === EMPTY));
  //   }
  // }, [boards]);

  useEffect(() => {
    if (!boards) {
      // get the board from the golang app
      GetBoard().then((board) => {
        console.log("first board",board)
        setBoards([board]);
      });
    }
  }, []);

  const startSolveTransition = () => {
    if (isSolving) return;
    setSolving(true);
    setSolveData(undefined);
    Solve().then((result) => {
      if (result.Status !== Status.SUCCESS) {
        alert("jotain meni pieleen yritä myöhemmin uudellen");
      } else {
        setSolving(false);
        setAnimating(true);
        setSolveData(result);
        setBoards(result.Iterations.map((item: any) => item.Board))
      }
    });
  };

  return (
    <ChakraProvider>
      <div id="App">
        {/* FIXME: remove these are for debugging  */}
        <Grid
          templateAreas={`"header header"
                  "nav main"
                  "nav footer"`}
          gridTemplateRows={"80px 1fr 30px"}
          gridTemplateColumns={"150px 1fr"}
          h="200px"
          gap="1"
          bg="white"
          color="blackAlpha.700"
          fontWeight="bold"
        >
          <GridItem pl="2" bg="orange.300" area={"header"}>
            <Button
              padding={5}
              margin={7}
              onClick={() => {
                startSolveTransition();
              }}
            >
              Solve!
            </Button>
            <Button
              padding={5}
              margin={7}
              onClick={() => {
                GenerateBoard().then((board) => {
                  setBoards([board])
                  console.log("reset board", board)
                  setAnimating(false)
                });
              }}
            >
              Reset!
            </Button>
          </GridItem>
          <GridItem
            pl="2"
            bg="pink.300"
            area={"nav"}
            style={{ textAlign: "start" }}
          >
            <Heading as="h3" size="lg">
              Data
            </Heading>
            {solveData && (
              <>
                <p>Solved in: {solveData.TimeElapsed}ms</p>
                <p>Moves: {solveData.IterationCount}</p>
              </>
            )}
          </GridItem>
          <GridItem pl="2" area={"main"}>
            {boards && boards.length > 0 && <Puzzle boards={boards} isAnimating={isAnimating} />}
          </GridItem>
        </Grid>
      </div>
    </ChakraProvider>
  );
}

function hash(boards : number[][]) : string {
  const str = boards.map(board => {
    return board.map(val => {
      return String.fromCharCode(val)
    }).join("");
  }).join("")
  return encodeURI(str);
}

const Puzzle = ({ boards, isAnimating }: { boards: number[][], isAnimating: boolean }) => {
  const index = useRef(0);
  const [board,setBoard] = useState(boards[0]);

  useEffect(()=>{
    if(boards.length > 1){
      let interval = setInterval(
        () => {
          if(index.current < boards.length - 1){
            index.current += 1
            console.log(`set board with index ${index.current}`,boards[index.current])
            setBoard(boards[index.current])
          }else{
            clearInterval(interval)
          }
        },300
      )
  
      return () => clearInterval(interval)
    }else if(boards.length === 1){
      console.log("set board 0 ");
      setBoard(boards[0])
    }
  },[hash(boards)])

  // const [board, cycleBoards,setBoards] = useCycle(...boards);
  // const [counter,setCounter] = useState(0);

  // let interval : number;
  // useEffect(()=>{
  //   setBoards(boards)
  //   clearInterval(interval);
  //   cycleBoards(0)
  // },[hash(boards)])
  // console.log(`boards length ${boards.length}`)
  // useEffect(() => {
  //   if(!isSolved(board) && isAnimating){
  //     interval = setInterval(
  //       () => {
  //         cycleBoards(counter)
  //         setCounter(counter+1)
  //       },300
  //     )
  
  //     return () => clearInterval(interval)
  //   }
  // }, [isAnimating]);  

  return (
    <div
      style={{
        display: "grid",
        gridTemplateColumns: "auto auto auto auto",
        gridGap: 10,
      }}
    >
      {board && board.map((item) => (
        <motion.div
          style={{
            width: 75,
            height: 75,
            borderRadius: 20,
            backgroundColor: "lightblue",
          }}
          key={item}
          layout
          transition={{ type: "spring", stiffness: 350, damping: 25 }}
        >
          <p> {item} </p>
        </motion.div>
      ))}
    </div>
  );
};


/**
 * Straight from framer-motion source like a cowboy 
 * 
 * Cycles through a series of visual properties. Can be used to toggle between or cycle through animations. It works similar to `useState` in React. It is provided an initial array of possible states, and returns an array of two arguments.
 *
 * An index value can be passed to the returned `cycle` function to cycle to a specific index.
 *
 * ```jsx
 * import * as React from "react"
 * import { motion, useCycle } from "framer-motion"
 *
 * export const MyComponent = () => {
 *   const [x, cycleX] = useCycle(0, 50, 100)
 *
 *   return (
 *     <motion.div
 *       animate={{ x: x }}
 *       onTap={() => cycleX()}
 *      />
 *    )
 * }
 * ```
 *
 * @param items - items to cycle through
 * @returns [currentState, cycleState]
 *
 * @public
 */
export function useCycle<T>(...propItems: T[]): [T, Cycle, React.Dispatch<React.SetStateAction<T[]>>] {
  const index = useRef(0)
  const [items, setItems] = useState(propItems);
  const [item, setItem] = useState(items[index.current])

  const runCycle = useCallback(
      (next?: number) => {
          index.current =
              typeof next !== "number"
                  ? wrap(0, items.length, index.current + 1)
                  : next
          // console.log("set item with index ", index.current)
          // console.log("items length", items.length)
          setItem(items[index.current])
      },
      // The array will change on each call, but by putting items.length at
      // the front of this array, we guarantee the dependency comparison will match up
      // eslint-disable-next-line react-hooks/exhaustive-deps
      [items.length, ...items]
  )
  useEffect(() => {
    console.log("new items length ", items.length)
    // @ts-ignore
  }, [hash(items)])
  
  console.log("return item", item)
  return [item, runCycle, setItems]
}

export const wrap = (min: number, max: number, v: number) => {
  const rangeSize = max - min
  return ((((v - min) % rangeSize) + rangeSize) % rangeSize) + min
}

export default App;

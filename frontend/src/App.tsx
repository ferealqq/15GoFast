import { useEffect, useState } from "react";
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
import { motion, useCycle } from "framer-motion";

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
            {boards && <Puzzle boards={boards} isAnimating={isAnimating} />}
          </GridItem>
        </Grid>
      </div>
    </ChakraProvider>
  );
}

const Puzzle = ({ boards, isAnimating }: { boards: number[][], isAnimating: boolean }) => {
  console.log(...boards)
  const [board, cycleBoards] = useCycle(...boards);
  const [counter,setCounter] = useState(0);

  useEffect(()=>{
    cycleBoards(0)
  },[boards])
  useEffect(() => {
    // if(!board) return;
    if(!isSolved(board)){
      let interval = setInterval(
        () => {
          // console.log(counter)
          cycleBoards(counter)
          setCounter(counter+1)
        },300
      )

      return () => clearInterval(interval)
    }
  }, [board, cycleBoards,isAnimating]);  

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

export default App;

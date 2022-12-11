import { useEffect, useRef, useState } from "react";
import "./App.css";
import { GenerateBoard, GetBoard, Solve } from "../wailsjs/go/main/App";
import {
  Button,
  ChakraProvider,
  Grid,
  GridItem,
  Heading,
} from "@chakra-ui/react";
import { motion } from "framer-motion";

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
  const [isSolving, setSolving] = useState(false);
  // Depricated
  const [isAnimating, setAnimating] = useState(false);
  const [solveData, setSolveData] = useState<SolveData | undefined>();
  // boards contains all the iterations of the board if it has been solved. If not this variable will only contain the starting state of the board 
  const [boards, setBoards] = useState<number[][] | undefined>();

  // set the initial board 
  useEffect(() => {
    if (!boards) {
      // get the board from the golang app
      GetBoard().then((board) => {
        console.log("first board",board)
        setBoards([board]);
      });
    }
  }, []);
  // solve the current board
  const startSolveTransition = () => {
    if (isSolving) return;
    setSolving(true);
    // clean the state of the application
    setSolveData(undefined);
    Solve().then((result) => {
      if (result.Status !== Status.SUCCESS) {
        alert("jotain meni pieleen yritä myöhemmin uudellen");
      } else {
        // start the solving board animation 
        setSolving(false);
        setAnimating(true);
        setSolveData(result);
        // Iterations contain all the different iterations of the board movements.
        setBoards(result.Iterations.map((item: any) => item.Board))
      }
    });
  };

  return (
    <ChakraProvider>
      <div id="App">
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
            setBoard(boards[index.current])
          }else{
            clearInterval(interval)
          }
        },300
      )
  
      return () => clearInterval(interval)
    }else if(boards.length === 1){
      setBoard(boards[0])
    }
  },[hash(boards)])

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

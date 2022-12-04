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

const swagStyle = {
  transition: "transform 0.3s ease 0s",
  animationDuration: "0.75s",
  transform: "translate3d(0px, 0px, 0px)",
  animationName: "swag",
};

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

enum Direction {
  DIRECTION_UP = 0,
  DIRECTION_DOWN = 1,
  DIRECTION_LEFT = 2,
  DIRECTION_RIGHT = 3,
}

type Move = {
  EmptyIndex: number;
  ToIndex: number;
  Direction: Direction;
};

type Cell = {
  value: number;
  move?: Move;
};

type Moves = { [cellValue: number]: Move | false };

const oppositeMove = (move: Move): Move => {
  switch (move.Direction) {
    case Direction.DIRECTION_UP:
      return {
        ...move,
        Direction: Direction.DIRECTION_DOWN,
      };
    case Direction.DIRECTION_DOWN:
      return {
        ...move,
        Direction: Direction.DIRECTION_UP,
      };

    case Direction.DIRECTION_LEFT:
      return {
        ...move,
        Direction: Direction.DIRECTION_RIGHT,
      };

    case Direction.DIRECTION_RIGHT:
      return {
        ...move,
        Direction: Direction.DIRECTION_LEFT,
      };

    default:
      return move;
  }
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
  const [puzzle, setPuzzle] = useState<undefined | Cell[]>();
  const [emptyIndex, setEmptyIndex] = useState<number | undefined>(undefined);
  const [isSolving, setSolving] = useState(false);
  const [isAnimating, setAnimating] = useState(false);
  const [solveData, setSolveData] = useState<SolveData | undefined>();
  const [boards, setBoards] = useState<number[][] | undefined>();

  useEffect(() => {
    if (puzzle) {
      setEmptyIndex(puzzle.findIndex((i) => i.value === EMPTY));
    }
  }, [puzzle]);

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
        let count = 0;
        // @ts-ignore
        setBoards(result.Iterations.map(item => item.Board))
        // const interval = setInterval(() => {
        //   if (count < result.Iterations.length) {
        //     const { Board, Move } = result.Iterations[count];
        //     const cells: Cell[] = Board.map((value: number, index: number) => {
        //       // if (index === Move.EmptyIndex || index === Move.ToIndex) {
        //       if (index === Move.EmptyIndex) {
        //         return { value, move: Move };
        //         // } else {
        //         //   return { value, move: oppositeMove(Move) };
        //         // }
        //       } else {
        //         return { value };
        //       }
        //     });
        //     setPuzzle(cells);
        //     count++;
        //   } else {
        //     setAnimating(false);
        //     clearInterval(interval);
        //   }
        // }, 800);
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
            {/* <Grid templateColumns="repeat(4, 4fr)" gap={6}>
              {!puzzle && <p>loading</p>}
              {puzzle &&
                puzzle.map(({ value, move }, index) => {
                  return (
                    <motion.div
                      animate={!!move ? "moving" : "stale"}
                      variants={createVariants(move)}
                      key={index}
                    >
                      <GridItem
                        w="100%"
                        h="100"
                        bg="blue.500"
                        style={swagStyle}
                      >
                        <Text>{move && "Moving"}</Text>
                        <Text fontSize="4xl">{value !== EMPTY && value}</Text>
                      </GridItem>
                    </motion.div>
                  );
                })}
            </Grid> */}
          </GridItem>
        </Grid>
      </div>
    </ChakraProvider>
  );
}

const Puzzle = ({ boards, isAnimating }: { boards: number[][], isAnimating: boolean }) => {
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
          console.log(counter)
          cycleBoards(counter)
          setCounter(counter+1)
        },300
      )

      return () => clearInterval(interval)
      // setTimeout(() => cycleBoards(), 300);
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

const createVariants = (move: Move | undefined) => {
  const moving: any = {
    // scale: [1, 2, 2, 1, 1],
    // rotate: [0, 0, 270, 270, 0],
    // borderRadius: ["20%", "20%", "50%", "50%", "20%"],
    duration: 0.8,
  };
  const stale = {
    scale: 0.9,
  };

  if (move === undefined) {
    return { stale, moving };
  }
  console.log(move);
  switch (move.Direction) {
    case Direction.DIRECTION_UP:
      moving["y"] = "calc(8vh + 10%)";
      break;

    case Direction.DIRECTION_DOWN:
      moving["y"] = "calc(8vh - 5%)";
      break;
    case Direction.DIRECTION_LEFT:
      moving["x"] = "calc(8vw + 5%)";
      break;

    case Direction.DIRECTION_RIGHT:
      moving["x"] = "calc(8vw - 5%)";
      break;

    default:
      break;
  }
  return {
    moving,
    stale,
  };
};

export default App;

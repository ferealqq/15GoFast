import { useEffect, useState } from "react";
import "./App.css";
import { GetBoard, Solve } from "../wailsjs/go/main/App";
import { Text, Button, ChakraProvider, Grid, GridItem } from "@chakra-ui/react";

const swagStyle = {
  transition: "transform 0.3s ease 0s",
  animationDuration: "0.75s",
  transform: "translate3d(0px, 0px, 0px)",
  animationName: "swag"
};

const EMPTY = 0

enum Status {
  FAILURE = 0,
  SUCCESS = 1,
  CUTOFF = 2,
  TIME_EXCEEDED = 3
}

function App() {
  const [puzzle, setPuzzle] = useState<undefined | number[]>();
  const [emptyIndex, setEmptyIndex] = useState<number | undefined>(undefined)

  useEffect(() => {
    if (puzzle) {
      setEmptyIndex(
        puzzle.findIndex(i => i === EMPTY)
      )
    }
  },[puzzle])

  useEffect(()=>{
    if(!puzzle){
      // get the board from the golang app 
      GetBoard().then(board => setPuzzle(board))
    }
  },[]);

  const startSolveTransition = () => {
    Solve().then(result => {
      if (result.Status !== Status.SUCCESS) {
        alert("jotain meni pieleen yritä myöhemmin uudellen")
      }
      setPuzzle(result.Iterations[result.Iterations.length-1])
    })
  }
  // swap the place of two columns in the react state
  // const swap = (to: number, from: number) => {
  //   const zero = puzzle[to];
  //   puzzle[to] = puzzle[from];
  //   puzzle[from] = zero;
  //   setPuzzle([...puzzle]);
  // };

  return (
    <ChakraProvider>
      <div id="App">
        {/* FIXME: remove these are for debugging  */}
        <Text fontSize="5xl">Empty index {emptyIndex}</Text>
        <Button
          onClick={() => {
           startSolveTransition() 
          }}
        >
          Solve!
        </Button>
        <Grid templateColumns="repeat(4, 4fr)" gap={6}>
          {!puzzle && <p>loading</p>}
          {puzzle && puzzle.map((number) => {
            return (
              <GridItem
                w="100%"
                h="100"
                bg="blue.500"
                style={swagStyle}
              >
                <Text fontSize="4xl">{number !== EMPTY && number}</Text>
              </GridItem>
            );
          })}
        </Grid>
      </div>
    </ChakraProvider>
  );
}

export default App;

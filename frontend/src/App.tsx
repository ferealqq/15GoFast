import { useEffect, useState } from "react";
import "./App.css";
import { GetRandomMove } from "../wailsjs/go/main/App";
import { Text, Button, ChakraProvider, Grid, GridItem } from "@chakra-ui/react";

const EMPTY = "EMPTY";

const numbers = [...[...Array(15).keys()].map((x) => x + 1), EMPTY];

function shuffle<T>(A: T[]): T[] {
  for (let i = A.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    const temp = A[i];
    A[i] = A[j];
    A[j] = temp;
  }
  return A;
}

const swagStyle = {
  transition: "transform 0.3s ease 0s",
  animationDuration: "0.75s",
  transform: "translate3d(0px, 0px, 0px)",
  animationName: "swag"
};

function App() {
  const [puzzle, setPuzzle] = useState(shuffle(numbers));
  const [emptyIndex, setEmptyIndex] = useState(puzzle.findIndex(x => x === EMPTY))

  useEffect(() => {
    setEmptyIndex(
      puzzle.findIndex(i => i === EMPTY)
    )
  },[puzzle])
  
  // swap the place of two columns in the react state
  const swap = (to: number, from: number) => {
    const zero = puzzle[to];
    puzzle[to] = puzzle[from];
    puzzle[from] = zero;
    setPuzzle([...puzzle]);
  };

  return (
    <ChakraProvider>
      <div id="App">
        {/* FIXME: remove these are for debugging  */}
        <Text fontSize="5xl">Empty index {emptyIndex}</Text>
        <Button
          onClick={() => {
            // TODO Types
            // @ts-ignore
            GetRandomMove(emptyIndex).then(res => swap(...res))
          }}
        >
          Change bro!
        </Button>
        <Grid templateColumns="repeat(4, 4fr)" gap={6}>
          {puzzle.map((number) => {
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

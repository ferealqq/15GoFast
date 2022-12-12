import { useEffect, useRef, useState } from "react";
import "./App.css";
import {
  GenerateBoard,
  GetBoard,
  Solve,
  SetComplexity,
  SetMaxRuntime,
} from "../wailsjs/go/main/App";
import {
  Button,
  ChakraProvider,
  Container,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Grid,
  GridItem,
  Heading,
  NumberInput,
  NumberInputField,
  SimpleGrid,
  Colors,
} from "@chakra-ui/react";
import { Field, Form, Formik } from "formik";
import { motion } from "framer-motion";
import Lottie, { LottieRefCurrentProps } from "lottie-react";
import atomAnimation from "./loading-atom.json";
// import { Token } from '@chakra-ui/styled-system'

enum Status {
  FAILURE = 0,
  SUCCESS = 1,
  CUTOFF = 2,
  TIME_EXCEEDED = 3,
}

type AppState =
  | "SOLVED"
  | "TIME_EXCEEDED"
  | "SOLVING"
  | "ANIMATING"
  | "WAITING"
  | "FAIL";

type SolveData = {
  Status: Status;
  Iterations: number[][];
  IterationCount: number;
  //time elapsed in milliseconds
  TimeElapsed: number;
};

const DEFAULT_MOVE_ANIMATION_TIME = 200;
const DEFAULT_COMPLEXITY = 1000;
const DEFAULT_MAX_RUNTIME = 5500;

function App() {
  // Depricated
  const [solveData, setSolveData] = useState<SolveData | undefined>();
  // boards contains all the iterations of the board if it has been solved. If not this variable will only contain the starting state of the board
  const [boards, setBoards] = useState<number[][] | undefined>();
  const [complexity, _setComplexity] = useState<number | undefined>();
  const [maxRuntime, _setMaxruntime] = useState<number | undefined>();
  const [moveAnimationTime, setMoveAnimationTime] = useState<number>(
    DEFAULT_MOVE_ANIMATION_TIME
  );
  const [appState, setAppState] = useState<AppState>("WAITING");
  // set the initial board
  useEffect(() => {
    if (!boards) {
      // get the board from the golang app
      GetBoard().then((board) => {
        setBoards([board]);
      });
    }
  }, []);

  const resetBoard = () => {
    // golang backend will brake down if we generate a new board while solving the current one.
    if (appState === "SOLVING") return;
    GenerateBoard().then((board) => {
      setBoards([board]);
      setAppState("WAITING");
    });
  };
  const validateNumber = (value: number): undefined | string => {
    let error;
    if (value < 1) {
      error = "Number has to be greater than 1";
    }
    return error;
  };
  const setComplexity = (comp: number) => {
    // if validate complexity returns something else than undefiend it means that the validation didn't go through
    if (validateNumber(comp)) return;
    // sync the go backend and the react frontend
    SetComplexity(comp);
    _setComplexity(comp);
    resetBoard();
  };

  const setMaxruntime = (max: number) => {
    SetMaxRuntime(max);
    _setMaxruntime(max);
  };

  // solve the current board
  const startSolveTransition = () => {
    if (appState === "SOLVING") return;
    setAppState("SOLVING");
    // clean the state of the application
    setSolveData(undefined);
    Solve().then((result) => {
      switch (result.Status) {
        case Status.SUCCESS:
          // start the solving board animation
          setAppState("ANIMATING");
          setSolveData(result);
          // Iterations contain all the different iterations of the board movements.
          setBoards(result.Iterations.map((item: any) => item.Board));

          break;
        case Status.TIME_EXCEEDED:
          setAppState("TIME_EXCEEDED");
          break;

        case Status.FAILURE:
          setAppState("FAIL");
          break;

        default:
          break;
      }
    });
  };
  return (
    <ChakraProvider>
      <SimpleGrid
        id="App"
        columns={1}
        bg="cyan.50"
        minHeight={"100vh"}
        minWidth={800}
      >
        <Grid
          templateAreas={`"header header"
                  "nav main"
                  "nav footer"`}
          gridTemplateRows={"80px 1fr 30px"}
          gridTemplateColumns={"330px 1fr"}
          h="200px"
          gap="1"
          color="blackAlpha.700"
          fontWeight="bold"
        >
          <GridItem
            pl="2"
            area={"header"}
            display="flex"
            alignContent="space-between"
            padding={3}
          >
            <Heading>15 Puzzle solver:</Heading>
            {/* <AnimatingText></AnimatingText> */}
            {appState === "SOLVING" && <SolvingAnimation />}
            {appState === "SOLVED" && <SolvedText />}
            {appState === "ANIMATING" && <AnimatingText />}
            {appState === "TIME_EXCEEDED" && <TimeExceededText />}
            {appState === "FAIL" && <FailText />}
          </GridItem>
          <GridItem
            pl="2"
            bg="teal.100"
            borderRadius={10}
            area={"nav"}
            style={{ textAlign: "start" }}
            paddingX={3}
            marginX={2}
          >
            <Container
              display={"flex"}
              alignContent="space-between"
              paddingTop={3}
            >
              <Button
                padding={3}
                margin="auto"
                onClick={() => {
                  startSolveTransition();
                }}
                bg="cyan.700"
                color="cyan.50"
              >
                Solve!
              </Button>
              <Button
                padding={3}
                margin="auto"
                bg="cyan.700"
                color="cyan.50"
                onClick={() => {
                  if (appState === "SOLVING") {
                    alert(
                      "Can't generate new board whilst solving the current one."
                    );
                    return;
                  }
                  resetBoard();
                }}
              >
                New!
              </Button>
            </Container>
            <Heading as="h3" size="lg" paddingTop={3}>
              Configurations
            </Heading>
            <Formik
              initialValues={{
                complexity: DEFAULT_COMPLEXITY,
                maxRuntime: DEFAULT_MAX_RUNTIME,
                moveAnimationTime: DEFAULT_MOVE_ANIMATION_TIME,
              }}
              onSubmit={(values, actions) => {}}
            >
              {(props: any) => (
                // TODO validation errors won't work because i didn't have time to configure
                <Form
                  style={{
                    paddingBottom: 15,
                    paddingLeft: 10,
                    paddingRight: 10,
                  }}
                >
                  <Field name="complexity" validate={validateNumber}>
                    {({ field, form }: { field: any; form: any }) => (
                      <FormControl
                        isInvalid={
                          form.errors.complexity && form.touched.complexity
                        }
                      >
                        <FormLabel>Complexity</FormLabel>
                        <NumberInput min={1} defaultValue={DEFAULT_COMPLEXITY}>
                          <NumberInputField
                            onChange={(evt) => {
                              // FIXME validateField so that the error message could be shown
                              // props.validateField("complexity");
                              setComplexity(parseInt(evt.target.value));
                            }}
                            value={complexity}
                            bg="cyan.50"
                          />
                        </NumberInput>
                        <FormErrorMessage>
                          {form.errors.complexity}
                        </FormErrorMessage>
                      </FormControl>
                    )}
                  </Field>
                  <Field name="maxRuntime" validate={validateNumber}>
                    {({ field, form }: { field: any; form: any }) => (
                      <FormControl
                        isInvalid={
                          form.errors.maxRuntime && form.touched.maxRuntime
                        }
                      >
                        <FormLabel>Time limit (ms)</FormLabel>
                        <NumberInput min={1} defaultValue={DEFAULT_MAX_RUNTIME}>
                          <NumberInputField
                            onChange={(evt) =>
                              setMaxruntime(parseInt(evt.target.value))
                            }
                            value={maxRuntime}
                            bg="teal.50"
                          />
                        </NumberInput>
                        <FormErrorMessage>
                          {form.errors.maxRuntime}
                        </FormErrorMessage>
                      </FormControl>
                    )}
                  </Field>

                  <Field
                    name="moveAnimation"
                    validate={(num: number) =>
                      num < 26
                        ? "Move animation has to be greater than 26"
                        : undefined
                    }
                  >
                    {({ field, form }: { field: any; form: any }) => (
                      <FormControl
                        isInvalid={
                          form.errors.moveAnimation &&
                          form.touched.moveAnimation
                        }
                      >
                        <FormLabel>Single move animation time (ms)</FormLabel>
                        <NumberInput
                          min={26}
                          defaultValue={DEFAULT_MOVE_ANIMATION_TIME}
                        >
                          <NumberInputField
                            onChange={(evt) => {
                              const val = parseInt(evt.target.value);
                              if (val > 26) {
                                setMoveAnimationTime(
                                  parseInt(evt.target.value)
                                );
                              }
                            }}
                            value={moveAnimationTime}
                            bg="teal.50"
                          />
                        </NumberInput>
                        <FormErrorMessage>
                          {form.errors.moveAnimation}
                        </FormErrorMessage>
                      </FormControl>
                    )}
                  </Field>
                </Form>
              )}
            </Formik>

            <Heading as="h3" size="lg">
              Data
            </Heading>
            <Container>
              {solveData && (
                <>
                  <p>Solved in: {solveData.TimeElapsed}ms</p>
                  <p>Moves: {solveData.IterationCount}</p>
                </>
              )}
            </Container>
          </GridItem>
          <GridItem pl="2" area={"main"}>
            {boards && boards.length > 0 && (
              <Puzzle
                boards={boards}
                isAnimating={appState === "ANIMATING"}
                moveAnimationTime={moveAnimationTime}
                onAnimationEnd={() => setAppState("SOLVED")}
              />
            )}
          </GridItem>
        </Grid>
      </SimpleGrid>
    </ChakraProvider>
  );
}

function hash(boards: number[][]): string {
  const str = boards
    .map((board) => {
      return board
        .map((val) => {
          return String.fromCharCode(val);
        })
        .join("");
    })
    .join("");
  return encodeURI(str);
}

const Puzzle = ({
  boards,
  isAnimating,
  moveAnimationTime,
  onAnimationEnd,
}: {
  boards: number[][];
  isAnimating: boolean;
  moveAnimationTime: number;
  onAnimationEnd: () => void;
}) => {
  const index = useRef(0);
  const [board, setBoard] = useState(boards[0]);

  useEffect(() => {
    if (boards.length > 1) {
      // Cycle through the boards one by one to create the illusion of moving the tiles
      let interval = setInterval(() => {
        if (index.current < boards.length - 1) {
          index.current += 1;
          setBoard(boards[index.current]);
        } else {
          onAnimationEnd();
          clearInterval(interval);
        }
      }, moveAnimationTime);

      return () => clearInterval(interval);
    } else if (boards.length === 1) {
      // reset the intervals index so that the next animations would be clean
      index.current = 0;
      setBoard(boards[0]);
    }
  }, [hash(boards)]);

  return (
    <div
      style={{
        display: "grid",
        gridTemplateColumns: "auto auto auto auto",
        gridGap: 30,
        padding: 10,
        paddingLeft: 20,
        paddingRight: 20,
      }}
    >
      {board &&
        board.map((item) => (
          <motion.div
            style={{
              width: "auto",
              height: 80,
              borderRadius: 15,
              backgroundColor: "lightblue",
              display: "flex",
              justifyContent: "center",
            }}
            key={item}
            layout
            transition={{
              type: "spring",
              stiffness: 350,
              damping: 20,
              duration: moveAnimationTime - 25,
            }}
          >
            <Heading textAlign={"center"} margin="auto">
              {" "}
              {item}{" "}
            </Heading>
          </motion.div>
        ))}
    </div>
  );
};
const containerVariants = {
  initial: {
    transition: {
      staggerChildren: 0.2,
    },
  },
  animate: {
    transition: {
      staggerChildren: 0.2,
    },
  },
};

const spanStyle = {
  fontSize: "var(--chakra-fontSizes-4xl)",
  lineHeight: "1.2",
  color: "var(--chakra-colors-cyan-800)",
  display: "block",
};

const textTransition = {
  duration: 0.9,
  yoyo: Infinity,
  ease: "easeInOut",
};

const FunkyText = ({ children }: { children: string }) => {
  const splitted = children.split("");
  return (
    <motion.div
      style={{ display: "flex" }}
      variants={containerVariants}
      initial="initial"
      animate="animate"
    >
      {splitted.map((char, i) => {
        return (
          <motion.span
            style={spanStyle}
            variants={{
              initial: {
                color: "#0BC5EA",
              },
              animate: {
                color: "#065666",
              },
            }}
            transition={textTransition}
            key={i}
          >
            {char}
          </motion.span>
        );
      })}
    </motion.div>
  );
};

const AnimatingText = () => {
  return (
    <Container display="flex" justifyContent={"center"}>
      <StatusText>Solved:</StatusText>
      <span style={{ paddingLeft: "10px" }} />
      <FunkyText> Animating</FunkyText>
    </Container>
  );
};

const SolvingAnimation = () => {
  const lot = useRef<LottieRefCurrentProps | null>(null);
  useEffect(() => {
    if (lot.current) {
      // @ts-ignore
      lot.current.setSpeed(3);
    }
  }, [lot.current]);

  return (
    <Container>
      <StatusText>Solving</StatusText>
      <Lottie
        lottieRef={lot}
        animationData={atomAnimation}
        style={{ paddingBottom: 30 }}
      />
    </Container>
  );
};

const StatusText = ({
  children,
  color = "cyan.800",
}: {
  children: string;
  color?: string;
}) => {
  return (
    <Heading as="h5" size="xl" color={color}>
      {children}
    </Heading>
  );
};

const StatusContainer = ({
  children,
  color = "cyan.800",
}: {
  children: string;
  color?: string;
}) => (
  <Container>
    <StatusText color={color}>{children}</StatusText>
  </Container>
);

const SolvedText = () => {
  return <StatusContainer>Solved!</StatusContainer>;
};

const TimeExceededText = () => (
  <StatusContainer color="red.400">
    Time limit exceeded, trying again with higher time limit!
  </StatusContainer>
);
const FailText = () => (
  <StatusContainer color="red.400">
    Unexpected error, maybe some funky math stuff!
  </StatusContainer>
);

export default App;

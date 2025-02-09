module Main exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onInput, onClick)
import Parser exposing (..)

-- Main
main =
    Browser.sandbox { init = init, update = update, view = view }

-- Modèle
type alias Model =
    { text : String
    , erreurs : List DeadEnd
    , result : List Instruction
    }

-- Définition des types d'instruction
type Instruction
    = Forward Int
    | Left Int
    | Right Int
    | Repeat Int (List Instruction)

init : Model
init =
    { text = ""
    , erreurs = []
    , result = []
    }

-- Message
type Msg
    = In String
    | ParseText

update : Msg -> Model -> Model
update msg model =
    case msg of
        In texte ->
          {model | text = texte}
        ParseText ->
            case read model.text of
                Ok instructions ->
                    { model | erreurs = [], result = instructions } -- Pas d'erreur, on met à jour result

                Err parsingErrors ->
                    { model | erreurs = parsingErrors, result = [] } -- Erreur, on met à jour erreurs et on vide result


-- Vue
view : Model -> Html Msg
view model =
    div []
        [ viewInput "text" "exemple: [Reapeat 100 [Forward 10]]: " model.text In
        , button [ onClick ParseText ] [ text "Lire" ]
        , div [] [ text "Résultat: ", text (instructionsToString model.result) ]
        ]

viewInput : String -> String -> String -> (String -> msg) -> Html msg
viewInput t p v toMsg =
  input [ type_ t, placeholder p, value v, onInput toMsg ] []

-- Convert a list of instructions to a string
instructionsToString : List Instruction -> String
instructionsToString instructions =
    let
        instructionToString instr =
            case instr of
                Forward n ->
                    "Forward " ++ String.fromInt n
                Left n ->
                    "Left " ++ String.fromInt n
                Right n ->
                    "Right " ++ String.fromInt n
                Repeat n instrs ->
                    "Repeat " ++ String.fromInt n ++ instructionsToString instrs
    in
    "[" ++ String.join ", " (List.map instructionToString instructions) ++ "]"

-- read

read : String -> Result (List DeadEnd) (List Instruction)
read input =
    run parseInstructions input


parseInstruction : Parser Instruction
parseInstruction =
    oneOf
        [ succeed Forward
            |. keyword "Forward"
            |. spaces
            |= int
        , succeed Left
            |. keyword "Left"
            |. spaces
            |= int
        , succeed Right
            |. keyword "Right"
            |. spaces
            |= int
        , succeed Repeat
            |. keyword "Repeat"
            |. spaces
            |= int
            |. spaces
            |= lazy (\_ -> sequence
                { start = "["
                , separator = ","
                , end = "]"
                , spaces = spaces
                , item = parseInstruction
                , trailing = Forbidden
                })
        ]

parseInstructions : Parser (List Instruction)
parseInstructions =
    sequence
        { start = "["
        , separator = ","
        , end = "]"
        , spaces = spaces
        , item = parseInstruction
        , trailing = Forbidden
        }


parseLoop : List Instruction -> Parser (Step (List Instruction) (List Instruction))
parseLoop acc =
    oneOf
        [ parseInstruction
            |> andThen (\instr -> succeed (Loop (instr :: acc)))
        , succeed (Done (List.reverse acc))
        ]


module Main exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onInput, onClick)
import Parser exposing (..)
import Svg exposing (..)
import Svg.Attributes exposing (..) 


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
        , button [ onClick ParseText ] [ Html.text "Lire" ]
        , div [] [ Html.text "Résultat: ", Html.text (instructionsToString model.result) ]
        , div [] [ renderSVG model.result ]
        ]

viewInput : String -> String -> String -> (String -> msg) -> Html msg
viewInput t p v toMsg =
  input [ Html.Attributes.type_ t, placeholder p, value v, onInput toMsg ] []

-- Convert a list of instructions to a string
instructionsToString : List Instruction -> String
instructionsToString instructions =
    let
        instructionToString instru =
            case instru of
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
            |> andThen (\instru -> succeed (Loop (instru :: acc)))
        , succeed (Done (List.reverse acc))
        ]

-- Fonction pour rendre les instructions en SVG
renderSVG : List Instruction -> Html Msg
renderSVG instructions =
    let
        -- État initial : position (x, y) et angle de rotation
        initialState =
            { x = 200, y = 200, angle = 0 }

        -- Fonction pour appliquer une instruction et mettre à jour l'état
        applyInstruction : Instruction -> { x : Int, y : Int, angle : Int } -> ( { x : Int, y : Int, angle : Int }, List (Svg Msg) )
        applyInstruction instr state =
            case instr of
                Forward n ->
                    let
                        newX = state.x + round (toFloat n * cos (degrees (toFloat state.angle)))
                        newY = state.y - round (toFloat n * sin (degrees (toFloat state.angle)))
                        lineElement =
                            line
                                [ x1 (String.fromInt state.x)
                                , y1 (String.fromInt state.y)
                                , x2 (String.fromInt newX)
                                , y2 (String.fromInt newY)
                                , stroke "black"
                                , strokeWidth "2"
                                ]
                                []
                    in
                    ( { x = newX, y = newY, angle = state.angle }, [ lineElement ] )

                Left n ->
                    ( { state | angle = state.angle + n }, [] )

                Right n ->
                    ( { state | angle = state.angle - n }, [] )

                Repeat n instrs ->
                    let
                        (newState, elements) =
                            List.foldl
                                (\instru (accState, accElements) ->
                                    let
                                        (updatedState, newElements) = applyInstruction instru accState
                                    in
                                    (updatedState, accElements ++ newElements)
                                )
                                (state, [])
                                (List.concat (List.repeat n instrs))
                    in
                    (newState, elements)

        -- Appliquer toutes les instructions et collecter les éléments SVG
        (finalState, svgElements) =
            List.foldl
                (\instr (accState, accElements) ->
                    let
                        (updatedState, newElements) = applyInstruction instr accState
                    in
                    (updatedState, accElements ++ newElements)
                )
                (initialState, [])
                instructions
    in
    svg
        [ Svg.Attributes.width "400", Svg.Attributes.height "400", viewBox "0 0 400 400" ]
        svgElements
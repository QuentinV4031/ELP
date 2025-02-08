module Main exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onInput,onClick)



-- MAIN


main : Program () Model Msg
main =
  Browser.sandbox { init = init, update = update, view = view }



-- MODEL


type alias Model =
  { texte : String
  , programme : List String
  }


init : Model
init =
  Model "" []



-- UPDATE


type Msg
  = Prog String
  | Ready


update : Msg -> Model -> Model
update msg model =
  case msg of
    Prog texte ->
      { model | texte = texte }

    Ready ->
      parser model



-- VIEW


view : Model -> Html Msg
view model =
  div []
    [ viewInput "text" "exemple: [Reapeat 100 [Forward 10]]" model.texte Prog
    , button [ onClick Ready ] [ text "Click when ready!" ]
    ]


viewInput : String -> String -> String -> (String -> msg) -> Html msg
viewInput t p v toMsg =
  input [ type_ t, placeholder p, value v, onInput toMsg ] []


parser : Model -> Model
parser model =
  { model | programme = model.programme ++ [model.texte] }
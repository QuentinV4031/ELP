module Main exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import TcTurtle as Turtle exposing (Command(..))

-- MAIN

main =
    Browser.sandbox { init = init, update = update, view = view }
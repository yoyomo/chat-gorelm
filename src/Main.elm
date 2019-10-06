module Main exposing (Action(..), State, initiaState, main, reducer, view)

import Browser
import Browser.Events exposing (onKeyPress)
import Html exposing (Html, button, div, input, li, text, ul)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick, onInput)


main =
    Browser.sandbox { init = initiaState, update = reducer, view = view }



-- MODEL


type alias State =
    { inputText : String
    , messages : List String
    }


initiaState : State
initiaState =
    { inputText = ""
    , messages = []
    }



-- UPDATE


type Action
    = InputChange String
    | KeyPress String
    | SendMessage


reducer : Action -> State -> State
reducer action state =
    case action of
        InputChange value ->
            { state | inputText = value }

        SendMessage ->
            { state
                | messages = List.append state.messages [ state.inputText ]
                , inputText = ""
            }

        KeyPress e ->
            if e.value == "Enter" then
                { state
                    | messages = List.append state.messages [ state.inputText ]
                    , inputText = ""
                }

            else
                state



-- VIEW


view : State -> Html Action
view state =
    div []
        [ renderList state.messages
        , input
            [ placeholder "Type a message", value state.inputText, onInput InputChange, onKeyPress KeyPress ]
            []
        , button [ onClick SendMessage ] [ text "Send" ]
        ]


renderList : List String -> Html msg
renderList lst =
    ul []
        (List.map (\l -> li [] [ text l ]) lst)

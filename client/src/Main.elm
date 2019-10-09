module Main exposing (Action(..), State, initiaState, main, reducer, view)

import Browser
import Html exposing (Attribute, Html, button, div, input, li, text, ul)
import Html.Attributes exposing (..)
import Html.Events exposing (keyCode, on, onClick, onInput)
import Json.Decode


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



-- VIEW


view : State -> Html Action
view state =
    div []
        [ renderList state.messages
        , input
            [ placeholder "Type a message", value state.inputText, onInput InputChange, onEnter SendMessage ]
            []
        , button [ onClick SendMessage ] [ text "Send" ]
        ]


onEnter : Action -> Attribute Action
onEnter action =
    let
        isEnter code =
            if code == 13 then
                Json.Decode.succeed action

            else
                Json.Decode.fail "not ENTER"
    in
    on "keydown" (Json.Decode.andThen isEnter keyCode)


renderList : List String -> Html action
renderList lst =
    ul []
        (List.map (\l -> li [] [ text l ]) lst)

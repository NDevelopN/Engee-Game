package gamedummy

import (
	"os"
	"testing"

	"github.com/google/uuid"
)

var testRID = uuid.NewString()
var testUID = uuid.NewString()

var emptyGame = GameDummy{}

func TestMain(m *testing.M) {
	setupSuite()
	code := m.Run()
	cleanupSuite()
	os.Exit(code)
}

func TestCreateDefaultGame(t *testing.T) {
	game, err := CreateDefaultGame(testRID)
	if err != nil {
		t.Fatalf("TestCreateDefaultGame: %v", err)
	}

	t.Cleanup(func() { cleanUpTestGame(game) })
}

func TestStartGame(t *testing.T) {
	testGame := setupDefaultGame(t)

	expected := *testGame
	expected.Status = ACTIVE

	err := testGame.StartGame()
	if !gamesAreEqual(*testGame, expected) || err != nil {
		t.Fatalf(`TestStartGame(Valid) = %v, %v, want %v, nil`, testGame, err, expected)
	}
}

func TestStartGameEmptyGame(t *testing.T) {
	testGame := emptyGame
	expected := emptyGame

	err := testGame.StartGame()
	if !gamesAreEqual(testGame, expected) || err == nil {
		t.Fatalf(`TestStartGame(EmptyGame) = %v, %v, want %v, err`, testGame, err, expected)
	}
}

func TestStartGameInvalidStatus(t *testing.T) {
	testGame := setupDefaultGame(t)
	expected := *testGame

	testGame.Status = ACTIVE
	expected.Status = ACTIVE

	err := testGame.StartGame()
	if !gamesAreEqual(*testGame, expected) || err == nil {
		t.Fatalf(`TestStartGame(InvalidStatus) = %v, %v, want %v, err`, testGame, err, expected)
	}
}

func TestStartGameRESET(t *testing.T) {
	testGame := setupDefaultGame(t)
	expected := *testGame

	testGame.Status = RESET
	expected.Status = ACTIVE

	err := testGame.StartGame()
	if !gamesAreEqual(*testGame, expected) || err != nil {
		t.Fatalf(`TestStartGame(RESET) = %v, %v, want %v, nil`, testGame, err, expected)
	}
}

func TestPauseGame(t *testing.T) {
	testGame := setUpTestGame(t)

	expected := *testGame
	expected.Status = PAUSED

	err := testGame.PauseGame()
	if !gamesAreEqual(*testGame, expected) || err != nil {
		t.Fatalf(`TestPauseGame(Valid) = %v, %v, want %v, nil`, testGame, err, expected)
	}
}

func TestPauseGameEmptyGame(t *testing.T) {
	testGame := emptyGame
	expected := emptyGame

	err := testGame.PauseGame()
	if !gamesAreEqual(testGame, expected) || err == nil {
		t.Fatalf(`TestPauseGame(EmptyGame) = %v, %v, want %v, err`, testGame, err, expected)
	}
}

func TestPauseGameInvalidStatus(t *testing.T) {
	testGame := setupDefaultGame(t)
	expected := *testGame

	err := testGame.PauseGame()
	if !gamesAreEqual(*testGame, expected) || err == nil {
		t.Fatalf(`TestPauseGame(InvalidStatus) = %v, %v, want %v, err`, testGame, err, expected)
	}
}

func TestPauseGamePAUSED(t *testing.T) {
	testGame := setUpTestGame(t)
	expected := *testGame

	testGame.Status = PAUSED
	expected.Status = ACTIVE

	err := testGame.PauseGame()
	if !gamesAreEqual(*testGame, expected) || err != nil {
		t.Fatalf(`TestPauseGame(PAUSED) = %v, %v, want %v, nil`, testGame, err, expected)
	}
}

func TestResetGame(t *testing.T) {
	testGame := setUpTestGame(t)

	expected := *testGame
	expected.Status = RESET

	err := testGame.ResetGame()
	if !gamesAreEqual(*testGame, expected) || err != nil {
		t.Fatalf(`TestResetGame(Valid) = %v, %v, want %v, nil`, testGame, err, expected)
	}
}

func TestResetGameEmptyGame(t *testing.T) {
	testGame := emptyGame
	expected := emptyGame

	err := testGame.ResetGame()
	if !gamesAreEqual(testGame, expected) || err == nil {
		t.Fatalf(`TestResetGame(EmptyGame) = %v, %v, want %v, err`, testGame, err, emptyGame)
	}
}

func TestResetGameInvalidStatus(t *testing.T) {
	testGame := setupDefaultGame(t)
	expected := *testGame

	err := testGame.ResetGame()
	if !gamesAreEqual(*testGame, expected) || err == nil {
		t.Fatalf(`TestResetGame(InvalidStatus) = %v, %v, want %v, err`, testGame, err, expected)
	}
}

func TestResetGamePAUSED(t *testing.T) {
	testGame := setUpTestGame(t)
	expected := *testGame

	testGame.Status = PAUSED
	expected.Status = RESET

	err := testGame.ResetGame()
	if !gamesAreEqual(*testGame, expected) || err != nil {
		t.Fatalf(`TestResetGame(PAUSED) = %v, %v, want %v, nil`, testGame, err, expected)
	}
}

func TestEndGame(t *testing.T) {
	testGame := setUpTestGame(t)

	err := testGame.EndGame()
	if err != nil {
		t.Fatalf(`TestEndGame(Valid) = %v, want nil`, err)
	}
}

func TestEndGameEmptyGame(t *testing.T) {
	testGame := emptyGame

	err := testGame.EndGame()
	if err == nil {
		t.Fatalf(`TestEndGame(EmptyGame) = %v, want err`, err)
	}
}

func TestEndGamePAUSED(t *testing.T) {
	testGame := setUpTestGame(t)
	testGame.Status = PAUSED

	err := testGame.EndGame()
	if err != nil {
		t.Fatalf(`TestEndGame(PAUSED) = %v, want nil`, err)
	}
}

func setupSuite() {

}

func cleanupSuite() {

}

func updateListener(message []byte) error {

	return nil
}

func setupDefaultGame(t *testing.T) *GameDummy {
	game, err := CreateDefaultGame(testRID)
	if err != nil {
		panic(err)
	}

	game.AddListener(updateListener)

	t.Cleanup(func() { cleanUpTestGame(game) })

	return game
}

func setUpTestGame(t *testing.T) *GameDummy {

	testGame := setupDefaultGame(t)
	testGame.StartGame()

	return testGame
}

func cleanUpTestGame(game *GameDummy) {
	game.EndGame()
}

func gamesAreEqual(first GameDummy, second GameDummy) bool {
	if first.Status != second.Status {
		return false
	}

	if first.Rules != second.Rules {
		return false
	}

	if len(first.Players) != len(second.Players) {
		return false
	}

	for index, plr := range first.Players {
		if plr != second.Players[index] {
			return false
		}
	}

	return true
}

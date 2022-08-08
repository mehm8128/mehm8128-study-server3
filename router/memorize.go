package router

import (
	"math/rand"
	"mehm8128_study_server/model"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Memorize struct {
	Name string `json:"name" db:"name"`
}
type Word struct {
	Word   string `json:"word" db:"word"`
	WordJp string `json:"wordJp" db:"word_jp"`
}
type QuizResponse struct {
	Answer     model.WordResponse   `json:"answer" db:"answer"`
	Selections []model.WordResponse `json:"selection" db:"selection"`
}

func getMemorizes(c echo.Context) error {
	ctx := c.Request().Context()
	memorizes, err := model.GetMemorizes(ctx)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if memorizes == nil {
		return echo.NewHTTPError(http.StatusOK, []*model.MemorizeResponse{})
	}
	return echo.NewHTTPError(http.StatusOK, memorizes)
}

func postMemorize(c echo.Context) error {
	var memorize Memorize
	err := c.Bind(&memorize)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}
	ctx := c.Request().Context()
	//todo:権限チェック
	res, err := model.CreateMemorize(ctx, memorize.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return echo.NewHTTPError(http.StatusOK, res)
}

func getMemorize(c echo.Context) error {
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	ctx := c.Request().Context()
	memorize, err := model.GetMemorize(ctx, ID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if memorize == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "memorize not found")
	}
	words, err := model.GetWords(ctx, ID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if words == nil {
		words = []model.WordResponse{}
	}
	memorize.Words = words
	return echo.NewHTTPError(http.StatusOK, memorize)
}

func postWord(c echo.Context) error {
	var word Word
	err := c.Bind(&word)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	ctx := c.Request().Context()
	//todo:権限チェック
	res, err := model.AddWord(ctx, ID, word.Word, word.WordJp)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return echo.NewHTTPError(http.StatusOK, res)
}

func getQuiz(c echo.Context) error {
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	ctx := c.Request().Context()
	words, err := model.GetWords(ctx, ID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if words == nil || len(words) < 4 {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot make quiz")
	}
	//配列をシャッフル、words配列の中身は答えとなる単語
	for i := 0; i < len(words); i++ {
		j := rand.Intn(len(words))
		words[i], words[j] = words[j], words[i]
	}
	var quiz []*QuizResponse
	for i := 0; i < len(words); i++ {
		//答えと選択肢、選択肢どうしが重複しないように選択肢を作る
		var numbers []int
		for j := 0; j < 4; j++ {
			numbers = append(numbers, rand.Intn(len(words)))
			//todo:選択肢どうしで重複しないようにする
		}
		for j := 0; j < len(numbers); j++ {
			if numbers[j] == i {
				numbers[j] = (numbers[j] + 1) % len(words)
			}
		}
		quiz = append(quiz, &QuizResponse{
			Answer:     words[i],
			Selections: []model.WordResponse{words[numbers[0]], words[numbers[1]], words[numbers[2]], words[numbers[3]]},
		})
		//答えの選択肢だけすり替え
		answerNumber := rand.Intn(4)
		quiz[i].Selections[answerNumber] = words[i]
	}

	return echo.NewHTTPError(http.StatusOK, quiz)
}

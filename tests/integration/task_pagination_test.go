package integration

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dherik/ddd-golang-project/internal/app/api"
	"github.com/dherik/ddd-golang-project/tests/integration/setup"
)

const tasksURL = "http://localhost:3333/tasks"

type taskPageBody struct {
	Items   []api.TaskResponse `json:"items"`
	HasMore bool               `json:"hasMore"`
}

func (s *TaskTestSuite) requestTaskList(token string, params url.Values) (int, []byte) {
	s.T().Helper()

	req, err := http.NewRequest(http.MethodGet, tasksURL+"?"+params.Encode(), nil)
	s.Require().NoError(err)
	req.Header.Set("Authorization", "Bearer "+token)

	response, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	s.Require().NoError(err)

	return response.StatusCode, body
}

func (s *TaskTestSuite) fetchTaskPage(token string, params url.Values) taskPageBody {
	s.T().Helper()

	status, body := s.requestTaskList(token, params)
	s.Require().Equal(http.StatusOK, status, string(body))

	var page taskPageBody
	s.Require().NoError(json.Unmarshal(body, &page))
	return page
}

func wideDateRange() url.Values {
	return url.Values{
		"startDate": {"2020-01-01T00:00:00Z"},
		"endDate":   {"2030-01-01T00:00:00Z"},
	}
}

func (s *TaskTestSuite) TestPaginationReturnsNewestFirstPageByDefault() {
	token := setup.Login("admin", "some_password")

	base := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 1; i <= 25; i++ {
		setup.InsertTaskWithCreatedAt(1, fmt.Sprintf("Paged task %02d", i), base.Add(time.Duration(i)*time.Minute))
	}

	page := s.fetchTaskPage(token, wideDateRange())

	s.Len(page.Items, 20)
	s.True(page.HasMore)
	s.Equal("Paged task 25", page.Items[0].Description)
	s.Equal("Paged task 06", page.Items[19].Description)
}

func (s *TaskTestSuite) TestPaginationHonoursLimitAndOffset() {
	token := setup.Login("admin", "some_password")

	base := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 1; i <= 25; i++ {
		setup.InsertTaskWithCreatedAt(1, fmt.Sprintf("Paged task %02d", i), base.Add(time.Duration(i)*time.Minute))
	}

	params := wideDateRange()
	params.Set("limit", "5")
	params.Set("offset", "3")

	page := s.fetchTaskPage(token, params)

	s.Len(page.Items, 5)
	s.True(page.HasMore)
	s.Equal("Paged task 22", page.Items[0].Description)
	s.Equal("Paged task 18", page.Items[4].Description)
}

func (s *TaskTestSuite) TestPaginationWalksEveryTaskExactlyOnce() {
	token := setup.Login("jane_smith", "some_password")

	base := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	for i := 1; i <= 25; i++ {
		setup.InsertTaskWithCreatedAt(3, fmt.Sprintf("Walk task %02d", i), base.Add(time.Duration(i)*time.Minute))
	}

	seen := map[string]bool{}
	var descriptions []string
	lastPage := 0

	for offset := 0; offset <= 25; offset += 10 {
		params := wideDateRange()
		params.Set("limit", "10")
		params.Set("offset", fmt.Sprint(offset))

		page := s.fetchTaskPage(token, params)
		lastPage = offset
		for _, task := range page.Items {
			s.Falsef(seen[task.Description], "task %q was returned twice", task.Description)
			seen[task.Description] = true
			descriptions = append(descriptions, task.Description)
		}
		if !page.HasMore {
			break
		}
	}

	s.Len(descriptions, 25)
	s.Equal(20, lastPage)
	s.Equal("Walk task 25", descriptions[0])
	s.Equal("Walk task 01", descriptions[24])
}

func (s *TaskTestSuite) TestPaginationStopsOnAFullFinalPage() {
	token := setup.Login("jane_smith", "some_password")

	base := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	for i := 1; i <= 20; i++ {
		setup.InsertTaskWithCreatedAt(3, fmt.Sprintf("Full page task %02d", i), base.Add(time.Duration(i)*time.Minute))
	}

	params := wideDateRange()
	params.Set("limit", "20")

	page := s.fetchTaskPage(token, params)

	s.Len(page.Items, 20)
	s.False(page.HasMore)
}

func (s *TaskTestSuite) TestPaginationTreatsOffsetZeroLikeNoOffset() {
	token := setup.Login("jane_smith", "some_password")

	base := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	for i := 1; i <= 5; i++ {
		setup.InsertTaskWithCreatedAt(3, fmt.Sprintf("Zero offset task %02d", i), base.Add(time.Duration(i)*time.Minute))
	}

	withoutOffset := s.fetchTaskPage(token, wideDateRange())

	params := wideDateRange()
	params.Set("offset", "0")
	withOffset := s.fetchTaskPage(token, params)

	s.Equal(withoutOffset.Items, withOffset.Items)
	s.Equal(withoutOffset.HasMore, withOffset.HasMore)
}

func (s *TaskTestSuite) TestPaginationPastTheEndReturnsAnEmptyArray() {
	token := setup.Login("jane_smith", "some_password")

	base := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	setup.InsertTaskWithCreatedAt(3, "Only task", base)

	params := wideDateRange()
	params.Set("offset", "10")

	status, body := s.requestTaskList(token, params)
	s.Require().Equal(http.StatusOK, status, string(body))

	var page taskPageBody
	s.Require().NoError(json.Unmarshal(body, &page))
	s.Empty(page.Items)
	s.False(page.HasMore)

	var envelope map[string]json.RawMessage
	s.Require().NoError(json.Unmarshal(body, &envelope))
	s.Equal("[]", strings.TrimSpace(string(envelope["items"])))
}

func (s *TaskTestSuite) TestPaginationBreaksTimestampTiesByDescendingId() {
	token := setup.Login("jane_smith", "some_password")

	sameInstant := time.Date(2025, 3, 1, 12, 0, 0, 0, time.UTC)
	setup.InsertTaskWithCreatedAt(3, "Tie A", sameInstant)
	setup.InsertTaskWithCreatedAt(3, "Tie B", sameInstant)
	setup.InsertTaskWithCreatedAt(3, "Tie C", sameInstant)

	firstParams := wideDateRange()
	firstParams.Set("limit", "2")
	first := s.fetchTaskPage(token, firstParams)

	s.Require().Len(first.Items, 2)
	s.Equal("Tie C", first.Items[0].Description)
	s.Equal("Tie B", first.Items[1].Description)
	s.True(first.HasMore)

	secondParams := wideDateRange()
	secondParams.Set("limit", "2")
	secondParams.Set("offset", "2")
	second := s.fetchTaskPage(token, secondParams)

	s.Require().Len(second.Items, 1)
	s.Equal("Tie A", second.Items[0].Description)
	s.False(second.HasMore)
}

func (s *TaskTestSuite) TestPaginationRejectsInvalidPageParameters() {
	token := setup.Login("admin", "some_password")

	cases := []struct {
		name    string
		param   string
		value   string
		message string
	}{
		{name: "zero limit", param: "limit", value: "0", message: "invalid limit"},
		{name: "negative limit", param: "limit", value: "-1", message: "invalid limit"},
		{name: "non numeric limit", param: "limit", value: "abc", message: "invalid limit"},
		{name: "limit above maximum", param: "limit", value: "101", message: "invalid limit"},
		{name: "negative offset", param: "offset", value: "-1", message: "invalid offset"},
		{name: "non numeric offset", param: "offset", value: "abc", message: "invalid offset"},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			params := wideDateRange()
			params.Set(tc.param, tc.value)

			status, body := s.requestTaskList(token, params)

			s.Require().Equal(http.StatusBadRequest, status, string(body))

			var errorBody struct {
				Error struct {
					Code    int    `json:"code"`
					Message string `json:"message"`
				} `json:"error"`
			}
			s.Require().NoError(json.Unmarshal(body, &errorBody))
			s.Equal(http.StatusBadRequest, errorBody.Error.Code)
			s.Equal(tc.message, errorBody.Error.Message)
		})
	}
}

func (s *TaskTestSuite) TestPaginationOnlyReturnsTheCallersTasks() {
	adminToken := setup.Login("admin", "some_password")
	johnToken := setup.Login("john_doe", "some_password")

	base := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	for i := 1; i <= 5; i++ {
		setup.InsertTaskWithCreatedAt(1, fmt.Sprintf("Admin task %02d", i), base.Add(time.Duration(i)*time.Minute))
		setup.InsertTaskWithCreatedAt(2, fmt.Sprintf("John task %02d", i), base.Add(time.Duration(i)*time.Minute))
	}

	adminPage := s.fetchTaskPage(adminToken, wideDateRange())
	s.Require().NotEmpty(adminPage.Items)
	for _, task := range adminPage.Items {
		s.Equal("1", task.UserId)
	}

	params := wideDateRange()
	params.Set("offset", "3")
	johnPage := s.fetchTaskPage(johnToken, params)
	s.Require().NotEmpty(johnPage.Items)
	for _, task := range johnPage.Items {
		s.Equal("2", task.UserId)
	}
}

func (s *TaskTestSuite) TestPaginationInsertDriftAndEndDatePinning() {
	token := setup.Login("admin", "some_password")

	base := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	for i := 1; i <= 5; i++ {
		setup.InsertTaskWithCreatedAt(1, fmt.Sprintf("Drift task %02d", i), base.Add(time.Duration(i)*time.Minute))
	}

	firstParams := wideDateRange()
	firstParams.Set("limit", "2")
	firstPage := s.fetchTaskPage(token, firstParams)

	s.Require().Len(firstPage.Items, 2)
	s.Equal("Drift task 05", firstPage.Items[0].Description)
	s.Equal("Drift task 04", firstPage.Items[1].Description)

	// A task created between page requests shifts the list down by one.
	setup.InsertTaskWithCreatedAt(1, "Drift task 06", base.Add(6*time.Minute))

	driftedParams := wideDateRange()
	driftedParams.Set("limit", "2")
	driftedParams.Set("offset", "2")
	driftedPage := s.fetchTaskPage(token, driftedParams)

	s.Require().Len(driftedPage.Items, 2)
	s.Equal("Drift task 04", driftedPage.Items[0].Description)

	// Pinning endDate to the first page's upper bound freezes the walk.
	pinnedParams := url.Values{
		"startDate": {"2020-01-01T00:00:00Z"},
		"endDate":   {firstPage.Items[0].CreatedAt.Format(time.RFC3339Nano)},
		"limit":     {"2"},
		"offset":    {"2"},
	}
	pinnedPage := s.fetchTaskPage(token, pinnedParams)

	s.Require().Len(pinnedPage.Items, 2)
	s.Equal("Drift task 03", pinnedPage.Items[0].Description)
	s.Equal("Drift task 02", pinnedPage.Items[1].Description)
}

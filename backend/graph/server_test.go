package graph_test

import (
	"bytes"
	stdcmp "cmp"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"iter"
	"maps"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
	"time"

	"careco/backend/domain"
	"careco/backend/domain/mock"
	"careco/backend/graph/test"
	"careco/backend/tests"

	"github.com/99designs/gqlgen/graphql"
	"github.com/aereal/optional"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.uber.org/mock/gomock"
	"gopkg.in/yaml.v3"
)

func TestServer(t *testing.T) {
	t.Parallel()

	root := new(casesRoot)
	if err := yaml.NewDecoder(bytes.NewReader(caseDefs)).Decode(root); err != nil {
		t.Fatal(err)
	}
	var idx int
	for tc := range root.All() {
		t.Run(tc.caseName, func(t *testing.T) {
			t.Parallel()

			buf := new(bytes.Buffer)
			if err := json.NewEncoder(buf).Encode(tc.Request); err != nil {
				t.Fatal(err)
			}

			ctrl := gomock.NewController(t)
			h := test.BuildHandler(ctrl)
			if c := mockCalls[tc.caseName]; c != nil {
				if c.drivingRecordCommand != nil {
					c.drivingRecordCommand(h.DrivingRecordCommand)
				}
				if c.drivingRecordQuery != nil {
					c.drivingRecordQuery(h.DrivingRecordQuery)
				}
			}
			srv := httptest.NewServer(h)
			t.Cleanup(srv.Close)

			req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, srv.URL, buf)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("content-type", "application/json")

			gotResp, err := srv.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer gotResp.Body.Close()
			assertsResponse(t, tc.Response, gotResp)
		})
		idx++
	}
}

//go:embed tests.yml
var caseDefs []byte

type gqlResponse struct {
	Errors []map[string]any `json:"errors,omitempty"`
	Data   map[string]any   `json:"data"`
}

type casesRoot struct {
	Cases map[string]*communication `json:"cases"`
}

type communication struct {
	Request  *graphql.RawParams `json:"request"`
	Response *gqlResponse       `json:"response"`
}

type testCase struct {
	*communication

	caseName string
}

func (r *casesRoot) All() iter.Seq[*testCase] {
	return func(yield func(*testCase) bool) {
		for _, caseName := range slices.SortedStableFunc(maps.Keys(r.Cases), stdcmp.Compare) {
			tc := &testCase{
				caseName:      caseName,
				communication: r.Cases[caseName],
			}
			if !yield(tc) {
				return
			}
		}
	}
}

func assertsResponse(t *testing.T, want *gqlResponse, got *http.Response) {
	t.Helper()
	jv, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	mv := make(map[string]any)
	if err := json.Unmarshal(jv, &mv); err != nil {
		t.Fatal(err)
	}
	exp := &responseExpectation{
		Status: http.StatusOK,
		Body:   mv,
	}
	if diff := cmp.Diff(exp, transformResponse(got)); diff != "" {
		t.Errorf("response (-want, +got):\n%s", diff)
	}
}

type responseExpectation struct {
	Status int
	Body   map[string]any
}

func transformResponse(hr *http.Response) *responseExpectation {
	defer hr.Body.Close()
	m := make(map[string]any)
	_ = json.NewDecoder(hr.Body).Decode(&m)
	return &responseExpectation{
		Status: hr.StatusCode,
		Body:   m,
	}
}

type mocks struct {
	drivingRecordCommand func(m *mock.MockDrivingRecordCommand)
	drivingRecordQuery   func(m *mock.MockDrivingRecordQuery)
}

var mockCalls = map[string]*mocks{
	"mutation recordDrivingRecord/ok": {
		drivingRecordCommand: func(m *mock.MockDrivingRecordCommand) {
			want := &domain.DrivingRecord{
				OdometerValue: 10001,
				Date:          time.Date(2025, time.October, 3, 0, 0, 0, 0, time.UTC),
			}
			m.EXPECT().RecordDrivingRecord(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, record *domain.DrivingRecord) error {
					return tests.Diff(want, record, cmpopts.EquateApproxTime(time.Millisecond), tests.EquateOptional[string]())
				}).
				Times(1)
		},
	},
	"query monthlyReport/ok": {
		drivingRecordQuery: func(m *mock.MockDrivingRecordQuery) {
			wantInterval := domain.Interval[time.Time]{
				Start: domain.ClosedEndpoint(time.Date(2025, time.October, 1, 0, 0, 0, 0, time.Local)),
				End:   domain.OpenEndpoint(time.Date(2025, time.November, 1, 0, 0, 0, 0, time.Local)),
			}
			ret := []*domain.DrivingRecord{
				{
					OdometerValue: 10003,
					Date:          time.Date(2025, time.October, 3, 0, 0, 0, 0, time.Local),
				},
				{
					OdometerValue: 10002,
					Date:          time.Date(2025, time.October, 2, 0, 0, 0, 0, time.Local),
				},
				{
					OdometerValue: 10001,
					Date:          time.Date(2025, time.October, 1, 0, 0, 0, 0, time.Local),
				},
			}
			m.EXPECT().
				FindRecordsInPeriod(gomock.Any(), eqTimeInterval(wantInterval), domain.OrderDirectionAsc, optional.None[int]()).
				Return(ret, nil).
				Times(1)
			m.EXPECT().
				FindLastRecordInPeriod(gomock.Any(), eqTimeInterval(wantInterval)).
				Return(ret[0], nil).
				Times(1)
		},
	},
	"query recentDrivingRecords/ok": {
		drivingRecordQuery: func(m *mock.MockDrivingRecordQuery) {
			ret := []*domain.DrivingRecord{
				{
					OdometerValue: 10004,
					Date:          time.Date(2025, time.October, 3, 12, 34, 56, 0, time.UTC),
				},
				{
					OdometerValue: 10003,
					Date:          time.Date(2025, time.October, 2, 12, 34, 56, 0, time.UTC),
					Memo:          optional.Some("blah blah"),
				},
				{
					OdometerValue: 10002,
					Date:          time.Date(2025, time.October, 1, 12, 34, 56, 0, time.UTC),
				},
				{
					OdometerValue: 10001,
					Date:          time.Date(2025, time.September, 30, 12, 34, 56, 0, time.UTC),
					Memo:          optional.Some("blah blah"),
				},
			}
			m.EXPECT().
				FindRecordsInPeriod(gomock.Any(), eqTimeInterval(domain.EmptyInterval[time.Time]()), domain.OrderDirectionDesc, cmpOptional(optional.Some(4))).
				Return(ret, nil).
				Times(1)
		},
	},
	"query totalStatistics/ok": {
		drivingRecordQuery: func(m *mock.MockDrivingRecordQuery) {
			m.EXPECT().
				FindLastRecordInPeriod(gomock.Any(), eqTimeInterval(domain.EmptyInterval[time.Time]())).
				Return(&domain.DrivingRecord{OdometerValue: 45678}, nil).
				Times(1)
		},
	},
	"query totalStatistics/empty": {
		drivingRecordQuery: func(m *mock.MockDrivingRecordQuery) {
			m.EXPECT().
				FindLastRecordInPeriod(gomock.Any(), eqTimeInterval(domain.EmptyInterval[time.Time]())).
				Return(nil, domain.ErrDrivingRecordNotFound).
				Times(1)
		},
	},
	"query yearlyReport/ok": {
		drivingRecordQuery: func(m *mock.MockDrivingRecordQuery) {
			octoberMonth := domain.Interval[time.Time]{
				Start: domain.ClosedEndpoint(time.Date(2025, time.October, 1, 0, 0, 0, 0, time.Local)),
				End:   domain.OpenEndpoint(time.Date(2025, time.November, 1, 0, 0, 0, 0, time.Local)),
			}
			septemberMonth := domain.Interval[time.Time]{
				Start: domain.ClosedEndpoint(time.Date(2025, time.September, 1, 0, 0, 0, 0, time.Local)),
				End:   domain.OpenEndpoint(time.Date(2025, time.October, 1, 0, 0, 0, 0, time.Local)),
			}
			yearInterval := domain.Interval[time.Time]{
				Start: domain.ClosedEndpoint(time.Date(2025, time.January, 1, 0, 0, 0, 0, time.Local)),
				End:   domain.OpenEndpoint(time.Date(2026, time.January, 1, 0, 0, 0, 0, time.Local)),
			}

			oct := []*domain.DrivingRecord{
				{
					OdometerValue: 10007,
					Date:          time.Date(2025, time.October, 3, 12, 34, 56, 0, time.UTC),
				},
				{
					OdometerValue: 10006,
					Date:          time.Date(2025, time.October, 2, 12, 34, 56, 0, time.UTC),
				},
				{
					OdometerValue: 10005,
					Date:          time.Date(2025, time.October, 1, 12, 34, 56, 0, time.UTC),
				},
			}
			sep := []*domain.DrivingRecord{
				{
					OdometerValue: 10004,
					Date:          time.Date(2025, time.September, 30, 12, 34, 56, 0, time.UTC),
				},
				{
					OdometerValue: 10003,
					Date:          time.Date(2025, time.September, 29, 12, 34, 56, 0, time.UTC),
				},
				{
					OdometerValue: 10002,
					Date:          time.Date(2025, time.September, 28, 12, 34, 56, 0, time.UTC),
				},
			}

			m.EXPECT().
				FindRecordsInPeriod(gomock.Any(), eqTimeInterval(yearInterval), domain.OrderDirectionAsc, optional.None[int]()).
				Return(slices.Concat(oct, sep), nil).
				Times(1)
			m.EXPECT().
				FindRecordsInPeriod(gomock.Any(), eqTimeInterval(octoberMonth), domain.OrderDirectionAsc, optional.None[int]()).
				Return(oct, nil).
				Times(1)
			m.EXPECT().
				FindRecordsInPeriod(gomock.Any(), eqTimeInterval(septemberMonth), domain.OrderDirectionAsc, optional.None[int]()).
				Return(sep, nil).
				Times(1)

			m.EXPECT().
				FindLastRecordInPeriod(gomock.Any(), eqTimeInterval(yearInterval)).
				Return(oct[0], nil).
				Times(1)
			m.EXPECT().
				FindLastRecordInPeriod(gomock.Any(), eqTimeInterval(octoberMonth)).
				Return(oct[0], nil).
				Times(1)
			m.EXPECT().
				FindLastRecordInPeriod(gomock.Any(), eqTimeInterval(septemberMonth)).
				Return(sep[0], nil).
				Times(1)
		},
	},
	"query tripDistance": {
		drivingRecordQuery: func(m *mock.MockDrivingRecordQuery) {
			octoberMonth := domain.Interval[time.Time]{
				Start: domain.ClosedEndpoint(time.Date(2025, time.October, 1, 0, 0, 0, 0, time.Local)),
				End:   domain.OpenEndpoint(time.Date(2025, time.November, 1, 0, 0, 0, 0, time.Local)),
			}
			septemberMonth := domain.Interval[time.Time]{
				Start: domain.ClosedEndpoint(time.Date(2025, time.September, 1, 0, 0, 0, 0, time.Local)),
				End:   domain.OpenEndpoint(time.Date(2025, time.October, 1, 0, 0, 0, 0, time.Local)),
			}
			yearInterval := domain.Interval[time.Time]{
				Start: domain.ClosedEndpoint(time.Date(2025, time.January, 1, 0, 0, 0, 0, time.Local)),
				End:   domain.OpenEndpoint(time.Date(2026, time.January, 1, 0, 0, 0, 0, time.Local)),
			}

			oct := []*domain.DrivingRecord{
				{
					OdometerValue: 10014,
					Date:          time.Date(2025, time.October, 3, 12, 34, 56, 0, time.Local),
				},
				{
					OdometerValue: 10009,
					Date:          time.Date(2025, time.October, 2, 12, 34, 56, 0, time.Local),
				},
				{
					OdometerValue: 10006,
					Date:          time.Date(2025, time.October, 1, 12, 34, 56, 0, time.Local),
				},
			}
			sep := []*domain.DrivingRecord{
				{
					OdometerValue: 10004,
					Date:          time.Date(2025, time.September, 30, 12, 34, 56, 0, time.Local),
				},
				{
					OdometerValue: 10003,
					Date:          time.Date(2025, time.September, 29, 12, 34, 56, 0, time.Local),
				},
				{
					OdometerValue: 10002,
					Date:          time.Date(2025, time.September, 28, 12, 34, 56, 0, time.Local),
				},
			}

			m.EXPECT().
				FindRecordsInPeriod(gomock.Any(), eqTimeInterval(yearInterval), domain.OrderDirectionAsc, optional.None[int]()).
				Return(slices.Concat(oct, sep), nil).
				Times(1)
			m.EXPECT().
				FindRecordsInPeriod(gomock.Any(), eqTimeInterval(octoberMonth), domain.OrderDirectionAsc, optional.None[int]()).
				Return(oct, nil).
				Times(1)
			m.EXPECT().
				FindRecordsInPeriod(gomock.Any(), eqTimeInterval(septemberMonth), domain.OrderDirectionAsc, optional.None[int]()).
				Return(sep, nil).
				Times(1)

			beforeOctober := domain.Interval[time.Time]{
				Start: domain.ClosedEndpoint(time.Date(2025, time.September, 30, 0, 0, 0, 0, time.Local)),
				End:   domain.OpenEndpoint(time.Date(2025, time.October, 1, 0, 0, 0, 0, time.Local)),
			}
			m.EXPECT().
				FindLastRecordInPeriod(gomock.Any(), eqTimeInterval(beforeOctober)).
				Return(sep[0] /* 9/30 */, nil).
				Times(1)

			beforeSeptember := domain.Interval[time.Time]{
				Start: domain.ClosedEndpoint(time.Date(2025, time.August, 1, 0, 0, 0, 0, time.Local)),
				End:   domain.OpenEndpoint(time.Date(2025, time.September, 1, 0, 0, 0, 0, time.Local)),
			}
			m.EXPECT().
				FindLastRecordInPeriod(gomock.Any(), eqTimeInterval(beforeSeptember)).
				Return(nil, domain.ErrDrivingRecordNotFound).
				Times(1)

			m.EXPECT().
				FindLastRecordInPeriod(gomock.Any(), eqTimeInterval(domain.Interval[time.Time]{
					Start: domain.ClosedEndpoint(time.Date(2025, time.September, 1, 0, 0, 0, 0, time.Local)),
					End:   domain.OpenEndpoint(time.Date(2025, time.October, 1, 0, 0, 0, 0, time.Local)),
				})).
				Return(&domain.DrivingRecord{OdometerValue: 10004}, nil).
				Times(1)
			m.EXPECT().
				FindLastRecordInPeriod(gomock.Any(), eqTimeInterval(domain.Interval[time.Time]{
					Start: domain.ClosedEndpoint(time.Date(2025, time.October, 1, 0, 0, 0, 0, time.Local)),
					End:   domain.OpenEndpoint(time.Date(2025, time.October, 2, 0, 0, 0, 0, time.Local)),
				})).
				Return(oct[2] /* 10/1 */, nil).
				Times(1)
			m.EXPECT().
				FindLastRecordInPeriod(gomock.Any(), eqTimeInterval(domain.Interval[time.Time]{
					Start: domain.ClosedEndpoint(time.Date(2025, time.October, 2, 0, 0, 0, 0, time.Local)),
					End:   domain.OpenEndpoint(time.Date(2025, time.October, 3, 0, 0, 0, 0, time.Local)),
				})).
				Return(oct[1] /* 10/2 */, nil).
				Times(1)

			m.EXPECT().
				FindLastRecordInPeriod(gomock.Any(), eqTimeInterval(domain.Interval[time.Time]{
					Start: domain.ClosedEndpoint(time.Date(2025, time.September, 27, 0, 0, 0, 0, time.Local)),
					End:   domain.OpenEndpoint(time.Date(2025, time.September, 28, 0, 0, 0, 0, time.Local)),
				})).
				Return(&domain.DrivingRecord{OdometerValue: 10001, Date: time.Date(2025, time.September, 27, 1, 0, 0, 0, time.Local)}, nil).
				Times(1)
			m.EXPECT().
				FindLastRecordInPeriod(gomock.Any(), eqTimeInterval(domain.Interval[time.Time]{
					Start: domain.ClosedEndpoint(time.Date(2025, time.September, 28, 0, 0, 0, 0, time.Local)),
					End:   domain.OpenEndpoint(time.Date(2025, time.September, 29, 0, 0, 0, 0, time.Local)),
				})).
				Return(sep[2] /* 9/28 */, nil).
				Times(1)
			m.EXPECT().
				FindLastRecordInPeriod(gomock.Any(), eqTimeInterval(domain.Interval[time.Time]{
					Start: domain.ClosedEndpoint(time.Date(2025, time.September, 29, 0, 0, 0, 0, time.Local)),
					End:   domain.OpenEndpoint(time.Date(2025, time.September, 30, 0, 0, 0, 0, time.Local)),
				})).
				Return(sep[1] /* 9/29 */, nil).
				Times(1)

			m.EXPECT().
				FindLastRecordInPeriod(gomock.Any(), eqTimeInterval(octoberMonth)).
				Return(oct[0], nil).
				Times(1)
			m.EXPECT().
				FindLastRecordInPeriod(gomock.Any(), eqTimeInterval(septemberMonth)).
				Return(sep[0], nil).
				Times(1)
		},
	},
}

func eqTimeInterval(want domain.Interval[time.Time]) gomock.Matcher {
	return &timeIntervalMatcher{Interval: want}
}

func cmpOptional[T comparable](want optional.Option[T]) gomock.Matcher {
	return gomock.WantFormatter(
		gomock.StringerFunc(func() string { return fmt.Sprintf("%#v", want) }),
		gomock.GotFormatterAdapter(
			gomock.GotFormatterFunc(formatDetailed),
			gomock.Cond(func(got optional.Option[T]) bool {
				return optional.Equal(want, got)
			}),
		),
	)
}

func formatDetailed(got any) string {
	return fmt.Sprintf("%#v", got)
}

type timeIntervalMatcher struct {
	domain.Interval[time.Time]
}

var (
	_ gomock.Matcher      = (*timeIntervalMatcher)(nil)
	_ gomock.GotFormatter = (*timeIntervalMatcher)(nil)
)

func (*timeIntervalMatcher) Got(got any) string {
	interval, ok := got.(domain.Interval[time.Time])
	if !ok {
		return fmt.Sprintf(" %#v", got)
	}
	return formatTimeInterval(interval)
}

func (m *timeIntervalMatcher) Matches(x any) bool {
	rhs, ok := x.(domain.Interval[time.Time])
	if !ok {
		return false
	}
	lhs := m.Interval
	return eqTimeEndpoint(lhs.Start, rhs.Start) && eqTimeEndpoint(lhs.End, rhs.End)
}

func (m *timeIntervalMatcher) String() string {
	return formatTimeInterval(m.Interval)
}

func formatEndpoint(endpoint domain.Endpoint[time.Time]) string {
	buf := new(bytes.Buffer)
	if endpoint.Open {
		buf.WriteString("  OPEN")
	} else {
		buf.WriteString("closed")
	}
	buf.WriteRune('(')
	buf.WriteString(endpoint.Value.Format(time.RFC3339Nano))
	buf.WriteRune(')')
	return buf.String()
}

func formatTimeInterval(interval domain.Interval[time.Time]) string {
	return fmt.Sprintf("{\n\tstart=%s\n\t  end=%s\n}", formatEndpoint(interval.Start), formatEndpoint(interval.End))
}

func eqTimeEndpoint(a, b domain.Endpoint[time.Time]) bool {
	return a.Open == b.Open && a.Value.Equal(b.Value)
}

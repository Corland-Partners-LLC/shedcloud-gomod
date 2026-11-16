package partnerapi

import (
	"context"
	"testing"
)

// TestRemainingRoutes pins method + path for the calls the resource tests do
// not assert on (plain gets and the order twins of the shared quote calls).
func TestRemainingRoutes(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		name         string
		method, path string
		call         func(c *Client) error
	}{
		{"leads get", "GET", "/partner/v1/leads/l1", func(c *Client) error { _, err := c.Leads.Get(ctx, "l1"); return err }},
		{"quotes get", "GET", "/partner/v1/quotes/q1", func(c *Client) error { _, err := c.Quotes.Get(ctx, "q1"); return err }},
		{"quotes cash terms set", "PUT", "/partner/v1/quotes/q1/cash-terms", func(c *Client) error {
			_, err := c.Quotes.SetCashTerms(ctx, "q1", CashTermsRequest{DownPayment: "1.00"})
			return err
		}},
		{"quotes rto preview", "POST", "/partner/v1/quotes/q1/rto-terms/preview", func(c *Client) error {
			_, err := c.Quotes.PreviewRTOTerms(ctx, "q1", RTOTermsComputeRequest{TermMonths: 36})
			return err
		}},
		{"quotes rto compute", "POST", "/partner/v1/quotes/q1/rto-terms/compute", func(c *Client) error {
			_, err := c.Quotes.ComputeRTOTerms(ctx, "q1", RTOTermsComputeRequest{TermMonths: 36})
			return err
		}},
		{"orders update", "PATCH", "/partner/v1/orders/o1", func(c *Client) error {
			d := "1.00"
			_, err := c.Orders.Update(ctx, "o1", OrderUpdateRequest{DiscountAmount: &d})
			return err
		}},
		{"orders customer", "PUT", "/partner/v1/orders/o1/customer", func(c *Client) error {
			_, err := c.Orders.SetCustomer(ctx, "o1", CustomerLinkRequest{Name: "Sam"})
			return err
		}},
		{"orders line build", "PATCH", "/partner/v1/orders/o1/lines/ln1/build", func(c *Client) error {
			_, err := c.Orders.UpdateLineBuild(ctx, "o1", "ln1", LineBuildRequest{ConfigurationID: "cfg"})
			return err
		}},
		{"orders cash preview", "POST", "/partner/v1/orders/o1/cash-terms/preview", func(c *Client) error {
			_, err := c.Orders.PreviewCashTerms(ctx, "o1", CashTermsComputeRequest{})
			return err
		}},
		{"orders cash compute", "POST", "/partner/v1/orders/o1/cash-terms/compute", func(c *Client) error {
			_, err := c.Orders.ComputeCashTerms(ctx, "o1", CashTermsComputeRequest{})
			return err
		}},
		{"orders rto preview", "POST", "/partner/v1/orders/o1/rto-terms/preview", func(c *Client) error {
			_, err := c.Orders.PreviewRTOTerms(ctx, "o1", RTOTermsComputeRequest{})
			return err
		}},
		{"orders rto compute", "POST", "/partner/v1/orders/o1/rto-terms/compute", func(c *Client) error {
			_, err := c.Orders.ComputeRTOTerms(ctx, "o1", RTOTermsComputeRequest{})
			return err
		}},
		{"orders rto set", "PUT", "/partner/v1/orders/o1/rto-terms", func(c *Client) error {
			_, err := c.Orders.SetRTOTerms(ctx, "o1", RTOTermsRequest{TermMonths: 36})
			return err
		}},
		{"orders mileage", "POST", "/partner/v1/orders/o1/mileage", func(c *Client) error { _, err := c.Orders.ComputeMileage(ctx, "o1"); return err }},
		{"orders tax", "POST", "/partner/v1/orders/o1/tax/resolve", func(c *Client) error { _, err := c.Orders.ResolveTax(ctx, "o1"); return err }},
		{"orders reprice preview", "POST", "/partner/v1/orders/o1/reprice/preview", func(c *Client) error {
			_, err := c.Orders.PreviewReprice(ctx, "o1")
			return err
		}},
		{"orders reprice", "POST", "/partner/v1/orders/o1/reprice", func(c *Client) error { _, err := c.Orders.Reprice(ctx, "o1"); return err }},
		{"locations get", "GET", "/partner/v1/locations/loc1", func(c *Client) error { _, err := c.Locations.Get(ctx, "loc1"); return err }},
		{"work orders get", "GET", "/partner/v1/work-orders/w1", func(c *Client) error { _, err := c.WorkOrders.Get(ctx, "w1"); return err }},
		{"users get", "GET", "/partner/v1/users/u1", func(c *Client) error { _, err := c.Users.Get(ctx, "u1"); return err }},
		{"users update", "PATCH", "/partner/v1/users/u1", func(c *Client) error {
			s := "inactive"
			_, err := c.Users.Update(ctx, "u1", UserUpdateRequest{Status: &s})
			return err
		}},
		{"agreements list", "GET", "/partner/v1/agreements", func(c *Client) error {
			_, err := c.Agreements.List(ctx, AgreementListParams{})
			return err
		}},
		{"agreements get", "GET", "/partner/v1/agreements/ag1", func(c *Client) error { _, err := c.Agreements.Get(ctx, "ag1"); return err }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, rec := fakeERP(t, 200, `{}`)
			if err := tc.call(c); err != nil {
				t.Fatal(err)
			}
			expect(t, rec, tc.method, tc.path)
		})
	}
}

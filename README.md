# concurrency

Gathered utilities to simplify concurrent programming, like for doing x number of request/response invocations in parallell.

## Install

```
go get github.com/dentech-floss/concurrency@v0.1.0
```

## Usage

### Concurrent Executor

With the [concurrent_executor.go](https://github.com/dentech-floss/concurrency/blob/main/pkg/concurrency/concurrent_executor.go) you can run x number of things in parallel, such as database lookups, with a timeout limiting the max amount of time allowed before completing. 

```go
package example

import (
    "github.com/dentech-floss/concurrency/pkg/concurrency"
)

type companyData struct {
    images      []*model.CompanyImage
    languages   []*model.CompanyLanguage
    memberships []*model.CompanyMembership
}

func (s *PatientGatewayServiceV1) fetchCompanyData(
    ctx context.Context,
    companyId int32,
) (*companyData, error) {

    executions := make([]*concurrency.Execution, 0)

    images := &concurrency.Execution{
        Request: func() (interface{}, error) {
            return s.repo.FindCompanyImagesByCompanyId(ctx, companyId)
        },
    }
    executions = append(executions, images)

    languages := &concurrency.Execution{
        Request: func() (interface{}, error) {
            return s.repo.FindCompanyLanguagesByCompanyId(ctx, companyId)
        },
    }
    executions = append(executions, languages)

    memberships := &concurrency.Execution{
        Request: func() (interface{}, error) {
            return s.repo.FindCompanyMembershipsByCompanyId(ctx, companyId)
        },
    }
    executions = append(executions, memberships)

    // Execute all of these lookup's in parallel, with a timeout
    err := concurrency.ExecuteConcurrently(executions, time.Duration(3)*time.Second)
    if err != nil {
        return nil, err
    }

    return &companyData{
        images:      images.Response.([]*model.CompanyImage),
        languages:   languages.Response.([]*model.CompanyLanguage),
        memberships: memberships.Response.([]*model.CompanyMembership),
    }, nil
}
```
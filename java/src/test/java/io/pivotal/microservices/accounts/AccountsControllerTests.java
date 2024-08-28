package io.pivotal.microservices.accounts;

import org.junit.Before;
import java.util.ArrayList;
import java.util.List;


/**
 * Test controller class.
 */
public class AccountsControllerTests extends AbstractAccountControllerTests {

    protected static final Account ACCOUNT_TEST = new Account(ACCOUNT_1, ACCOUNT_1_NAME);

    /**
     * Testing account repository.
     */
    protected static class TestAccountRepository implements AccountRepository {

        @Override
        public Account findByNumber(String accountNumber) {
            if (accountNumber.equals(ACCOUNT_1)) {
                return ACCOUNT_TEST;
            } else {
                return null;
            }
        }

        @Override
        public List<Account> findByOwnerContainingIgnoreCase(String partialName) {
            List<Account> accounts = new ArrayList<Account>();

            if (ACCOUNT_1_NAME.toLowerCase().indexOf(partialName.toLowerCase()) != -1) {
                accounts.add(ACCOUNT_TEST);
            }

            return accounts;
        }

        @Override
        public int countAccounts() {
            return 1;
        }
    }

    protected TestAccountRepository testRepo = new TestAccountRepository();

    @Before
    public void setup() {
        accountController = new AccountsController(testRepo);
    }
}

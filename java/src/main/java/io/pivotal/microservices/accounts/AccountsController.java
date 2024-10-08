package io.pivotal.microservices.accounts;

import io.pivotal.microservices.exceptions.AccountNotFoundException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import java.util.List;
import java.util.logging.Logger;

/**
 * A RESTFul controller for accessing account information.
 *
 * @author Paul Chapman
 */
@RestController
public class AccountsController {

    protected Logger logger = Logger.getLogger(AccountsController.class.getName());
    protected AccountRepository accountRepository;

    /**
     * Create an instance plugging in the respository of Accounts.
     * 
     * @param accountRepository An account repository implementation.
     */
    @Autowired
    public AccountsController(AccountRepository accountRepository) {
        this.accountRepository = accountRepository;

        logger.info("AccountRepository says system has " + accountRepository.countAccounts() + " accounts");
    }

    /**
     * Fetch an account with the specified account number.
     * 
     * @param accountNumber A numeric, 9 digit account number.
     * @return The account if found.
     */
    @RequestMapping("/accounts/{accountNumber}")
    public Account byNumber(@PathVariable("accountNumber") String accountNumber) {

        logger.info("accounts-service byNumber() invoked: " + accountNumber);
        Account account = accountRepository.findByNumber(accountNumber);
        logger.info("accounts-service byNumber() found: " + account);

        if (account == null) {
            throw new AccountNotFoundException(accountNumber);
        } else {
            return account;
        }
    }

    /**
     * Fetch accounts with the specified name. A partial case-insensitive match is
     * supported. So <code>http://.../accounts/owner/a</code> will find any accounts
     * with upper or lower case 'a' in their name.
     * 
     * @param partialName Name of the owener to fetch
     * @return A non-null, non-empty set of accounts.
     */
    @RequestMapping("/accounts/owner/{name}")
    public List<Account> byOwner(@PathVariable("name") String partialName) {
        logger.info("accounts-service byOwner() invoked: " + accountRepository.getClass().getName() + " for "
                + partialName);

        List<Account> accounts = accountRepository.findByOwnerContainingIgnoreCase(partialName);
        logger.info("accounts-service byOwner() found: " + accounts);

        if (accounts == null || accounts.size() == 0) {
            throw new AccountNotFoundException(partialName);
        } else {
            return accounts;
        }
    }
}

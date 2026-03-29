import logging
import time
from typing import Any, Callable
from io import BytesIO
from selenium import webdriver
from selenium.webdriver.common.by import By
from PIL import Image

LOGIN_QR_CODE_CSS_SELECTOR = "._akau"

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


class WhatsappDriver:
    def __init__(self, loginCallback: Callable[..., Any]):
        self.loginCallback = loginCallback
        self.driver = None

    def _check_and_create_driver(self):
        """Checks that Webdriver exists and is logged into WhatsApp
        Creates webdriver and calls _login() if necessary
        """
        # check that drive exists
        if self.driver is None:
            logger.info("Creating new Chrome webdriver")
            self.driver = webdriver.Chrome()

        # check that driver is on whatsapp
        if "WhatsApp" not in self.driver.title:
            logger.info("Changing website to https://web.whatsapp.com")
            self.driver.get("https://web.whatsapp.com")
            time.sleep(1)

        # check if login qr code is present -> login required
        if self._get_login_element():
            logger.info("Not logged in -> calling _login()")
            self._login()

    def _get_login_element(self):
        assert self.driver
        return self.driver.find_element(By.CSS_SELECTOR, LOGIN_QR_CODE_CSS_SELECTOR)

    def _login(self):
        """Performs login by extracting login qr code
        and calling the user provided loginCallback function
        """
        assert self.driver

        if not self._get_login_element():
            logger.warning("No login qr code found")
            return

        while "loading" in self._get_login_element().get_attribute("innerHTML"):
            time.sleep(0.2)

        time.sleep(0.2)
        self.driver.execute_script(
            "arguments[0].scrollIntoView();", self._get_login_element()
        )

        time.sleep(0.2)
        self._add_background_to_image(self._get_login_element().screenshot_as_png)

        time.sleep(20)

    def _add_background_to_image(self, png_bytes):
        print(type(png_bytes))
        img_bytes = BytesIO(png_bytes)
        print(type(img_bytes))
        img = Image.open(img_bytes).convert("RGB")

        border = 40
        new_w = img.width + 2 * border
        new_h = img.height + 2 * border

        background = Image.new("RGB", (new_w, new_h), (255, 255, 255))
        background.paste(img, (border, border))
        background.show()


whatsapp = WhatsappDriver(lambda: None)
whatsapp._check_and_create_driver()

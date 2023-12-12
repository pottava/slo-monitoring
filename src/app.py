import os

import streamlit as st
from streamlit.components.v1 import html


st.title("Streamlit sample")
rev = os.getenv("K_REVISION")
if rev:
    html(f'<script>console.log("Revision: {rev}");</script>', height=1)
    st.text(f"Revision: {rev}")
else:
    st.text("Hello!")

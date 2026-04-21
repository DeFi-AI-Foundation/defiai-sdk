"""Setup script for Patent N Python SDK"""

from setuptools import setup, find_packages

with open("README.md", "r", encoding="utf-8") as fh:
    long_description = fh.read()

setup(
    name="patent-n-sdk",
    version="1.0.0",
    author="Patent N Licensing",
    author_email="support@patent-n.example.com",
    description="Python SDK for Patent N Licensing API - OR_CCR_61 Bypass System",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/patent-n/sdk-python",
    packages=find_packages(),
    classifiers=[
        "Development Status :: 5 - Production/Stable",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
        "Programming Language :: Python :: 3.12",
    ],
    python_requires=">=3.8",
    install_requires=[
        "requests>=2.31.0",
        "urllib3>=2.0.0",
    ],
    extras_require={
        "dev": [
            "pytest>=7.4.0",
            "pytest-cov>=4.1.0",
            "black>=23.7.0",
            "mypy>=1.5.0",
        ],
    },
    keywords="patent-n or-ccr-61 bypass prepaid-card error-detection fintech",
    project_urls={
        "Bug Reports": "https://github.com/patent-n/sdk-python/issues",
        "Source": "https://github.com/patent-n/sdk-python",
        "Documentation": "https://docs.patent-n.example.com",
    },
)
